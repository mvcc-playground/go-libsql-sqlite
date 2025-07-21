De forma simples:

- **`_ "github.com/tursodatabase/libsql-client-go/libsql"`** (puro Go):

  - É um driver **sem cgo** (escrito apenas em Go).
  - Conecta ao Turso via **WebSocket/HTTP**, ideal para conexões remotas.
  - Não embute código C do SQLite, evitando conflitos de linker.
  - Leve, não precisa de compilador C ou bibliotecas extras.
  - Usado para acessar bancos remotos do Turso diretamente.

- **`_ "github.com/tursodatabase/go-libsql"`** (com cgo):
  - Usa **cgo** (embute código C do SQLite).
  - Suporta **réplicas embutidas locais** (sincroniza banco local com Turso).
  - Mais pesado, requer compilador C e pode causar conflitos com outros drivers SQLite (como `mattn/go-sqlite3`).
  - Usado quando você quer um banco local que sincroniza com o Turso.

**Resumo**: Use `libsql-client-go/libsql` para conexões remotas simples (como no seu caso). Use `go-libsql` se precisar de réplicas locais sincronizadas.
