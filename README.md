# xmlgen

XML generation with an emphasis on _explicit_ declaration of the resulting
structure.

`encoding/xml` is a mess compounded by code generators that run in to issues
such as [conflicting names](https://github.com/golang/go/issues/18564) that
can't reasonably be solved with the abstractions at hand.

The goal here is that the developer is in *complete* control of the XML
structure generated, but with a reasonably compact API.

For examples, see the cmd directory.

## TODO

* Namespaces. Yay...
