package parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"testing"
)

var q1 = `
{
	al(xid: "alice") {
		status
		_xid_
		follows {
			status
			_xid_
			follows {
				status
				_xid_
				follows {
					_xid_
					status
				}
			}
		}
		status
		_xid_
	}
}
`

var q2 = `query queryName {
		me(uid : "0x0a") {
			friends {
				name
			}
			gender,age
			hometown
		}
	}
`

var q3 = `
{
  debug(xid: "m.0bxtg") {
    type.object.name.en
    film.actor.film {
      film.performance.film {
        film.film.directed_by {
          type.object.name.en
        }
      }
    }
  }
}
`

var q4 = `
{
  debug(_xid_: "m.06pj8") {
    type.object.name.en
    film.director.film {
      type.object.name.en
      film.film.initial_release_date
      film.film.country
      film.film.starring {
        film.performance.actor {
          type.object.name.en
        }
        film.performance.character {
          type.object.name.en
        }
      }
      film.film.genre {
        type.object.name.en
      }
    }
  }
}
`

func TestQueryParse(t *testing.T) {
	input := antlr.NewInputStream(q4)
	lexer := NewGraphQLPMLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := NewGraphQLPMParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
	tree := p.Document()
	antlr.ParseTreeWalkerDefault.Walk(newMyListener(), tree)
}

func runParser(q string, b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := antlr.NewInputStream(q)
		lexer := NewGraphQLPMLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := NewGraphQLPMParser(stream)
		p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
		p.BuildParseTrees = true
		// uptill here we have a cost of : 19000 for q1
		// next call makes it 100 times more costly to : 1800000
		_ = p.Document()
	}
}

func BenchmarkQuery(b *testing.B) {
	b.Run("q1", func(b *testing.B) { runParser(q1, b) })
	b.Run("q2", func(b *testing.B) { runParser(q2, b) })
}
