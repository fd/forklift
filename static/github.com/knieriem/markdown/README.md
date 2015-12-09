This is an implementation of John Gruber's [markdown][] in
[Go][].  It is a translation of [peg-markdown][], written by
John MacFarlane in C, into Go.  It is using a modified version
of Andrew J Snodgrass' PEG parser [peg][] -- now supporting
LEG grammars --, which itself is based on the parser used
by peg-markdown.

[markdown]: http://daringfireball.net/projects/markdown/
[peg-markdown]: https://github.com/fd/forklift/static/github.com/jgm/peg-markdown
[peg]: https://github.com/fd/forklift/static/github.com/pointlander/peg
[Go]: http://golang.org/

Support for HTML and groff mm output is implemented, but LaTeX
output has not been ported. The output is identical
to that of peg-markdown.

I try to keep the grammar in sync with the C version, by
cherry-picking relevant changes. In the commit history the
corresponding revisions have a suffix *[jgm/peg-markdown].*

The Markdown parser has a performance similar to that of 
the original C version, and consumes less memory.

## Installation

Provided you have a copy of Go 1, and git is available,

	go get github.com/fd/forklift/static/github.com/knieriem/markdown

should download and install the package according to
your GOPATH settings.

See doc.go for an example how to use the package.

---

To create the command line program *markdown,* run

	go build github.com/fd/forklift/static/github.com/knieriem/markdown/cmd/markdown

the binary should then be available in the current directory.

To run tests, type

	go test github.com/fd/forklift/static/github.com/knieriem/markdown

At the moment, tests are based on the .text files from the
Markdown 1.0.3 test suite created by John Gruber, [imported from
peg-markdown][testsuite]. The output of the conversion of these
.text files to html is compared to the output of peg-markdown.

[testsuite]: https://github.com/fd/forklift/static/github.com/jgm/peg-markdown/tree/master/MarkdownTest_1.0.3

## Development

There is not yet a way to create a Go source file like
`parser.leg.go` automatically from another file, `parser.leg`,
when building packages and commands using the Go tool.  To make
*markdown* installable using `go get`, `parser.leg.go` has
been added to the VCS.

`Make parser` will update `parser.leg.go` using `leg` – which
is part of [knieriem/peg][] at github –, if parser.leg has
been changed, or if the Go file is missing. If a copy of *peg*
is not yet present on your system, run

	go get github.com/fd/forklift/static/github.com/knieriem/peg

Then `make parser` should succeed.

[knieriem/peg]: https://github.com/fd/forklift/static/github.com/knieriem/peg


## Extensions

In addition to the extensions already present in peg-markdown,
this package also supports definition lists (option `-dlists`)
similar to the way they are described in the documentation of
[PHP Markdown Extra][].

Definitions (`<dd>...</dd>`) are implemented using [`ListTight`][ListTight]
and `ListLoose`, on which bullet lists and ordered lists are based
already. If there is an empty line between the definition title and
the first definition, a loose list is expected, a tight list otherwise.

As definition item markers both `:` and `~` can be used.

[PHP Markdown Extra]: http://michelf.com/projects/php-markdown/extra/#def-list
[ListTight]: https://github.com/fd/forklift/static/github.com/knieriem/markdown/blob/master/parser.leg#L191


## Todo

*	Port tables and perhaps other extensions from [fletcher/peg-multimarkdown][mmd].

## Subdirectory Index

*	cmd/markdown	– command line program `markdown`

[mmd]: https://github.com/fd/forklift/static/github.com/fletcher/peg-multimarkdown
