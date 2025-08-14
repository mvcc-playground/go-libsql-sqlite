GROK resolveu

Sim, o erro de "duplicate symbols" ocorre porque você está usando o pacote `github.com/tursodatabase/go-libsql` (que requer cgo e embute o código C do SQLite) junto com o driver oficial do GORM para SQLite (`gorm.io/driver/sqlite`), que internamente importa `github.com/mattn/go-sqlite3` (também com cgo e SQLite embutido). Isso causa conflito de símbolos duplicados no linker, pois ambos incluem as mesmas funções do SQLite (como `_sqlite3_vsnprintf`, etc.).

O seu outro código funciona porque não usa o GORM, evitando o import do `mattn/go-sqlite3`.

### Solução para conexões remotas (como no seu código):

Seu código usa uma URL remota (`libsql://...`), não réplicas embutidas locais. Para isso, use o driver puro Go (sem cgo) do Turso: `github.com/tursodatabase/libsql-client-go/libsql`. Ele se comunica via protocolo WebSocket/HTTP e não embute código C, evitando o conflito.

1. **Atualize as dependências**:

   ```
   go get github.com/tursodatabase/libsql-client-go/libsql
   go mod tidy
   ```

   (Remova ou ignore o `go-libsql` se não precisar de réplicas embutidas.)

2. **Ajuste o código**:
   Mude apenas o import blank:

   ```go
   package main

   import (
   	"database/sql"
   	"fmt"
   	"os"

   	_ "github.com/tursodatabase/libsql-client-go/libsql"  // Mude para este (puro Go)
   	"gorm.io/driver/sqlite"
   	"gorm.io/gorm"
   )

   type Task struct {
   	gorm.Model
   	ID    int
   	Title int
   }

   func main() {
   	// dbName := "file:./local.db"
   	dbURL := fmt.Sprintf("libsql://%v-%v.aws-us-east-1.turso.io?authToken=%v",
   		env.DB_NAME,
   		env.DB_TURSO_USER,
   		env.DB_TOKEN,
   )

   	conn, err := sql.Open("libsql", dbHost)
   	if err != nil {
   		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
   		os.Exit(1)
   	}
   	defer conn.Close()
   	db, err := gorm.Open(sqlite.New(sqlite.Config{
   		Conn: conn,
   	}))
   	if err != nil {
   		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
   		os.Exit(1)
   	}
   	db.AutoMigrate(&Task{})
   }
   ```

3. **Execute**:
   ```
   go run main.go
   ```
   Isso deve compilar e rodar sem erros de linker, pois não há mais conflito de cgo/SQLite.

### Se você precisar de réplicas embutidas (com cgo):

- Use o fork `github.com/ytsruh/gorm-libsql` (mesmo arquivado, funciona e evita o import do `mattn/go-sqlite3`).
- Instale: `go get github.com/ytsruh/gorm-libsql`
- Ajuste o import: `import libsql "github.com/ytsruh/gorm-libsql"`
- Use o exemplo de "Full Feature Usage" do README (com `libsql.NewEmbeddedReplicaConnector` e `Conn: conn`).
- Mantenha `_ "github.com/tursodatabase/go-libsql"`.

Se ainda tiver problemas, verifique a versão do GORM/SQLite (atualize para a latest: `go get gorm.io/gorm@latest` e `go get gorm.io/driver/sqlite@latest`) ou compartilhe o output de `go env` para mais detalhes.
