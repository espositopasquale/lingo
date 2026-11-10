# Lessici italiani

Le raccolte in `tokenizer/italian_lexicon.go` sono state adattate il 5 novembre
2026 dai due file forniti localmente, senza dipendenze runtime da Moses o spaCy.
Le voci non portano metadati individuali.

- [Moses, nonbreaking_prefix.it](https://github.com/moses-smt/mosesdecoder/blob/master/scripts/share/nonbreaking_prefixes/nonbreaking_prefix.it):
  conversione dei prefissi in forme con punto, conservazione delle condizioni
  numeriche e gestione separata delle iniziali maiuscole. Moses è distribuito
  sotto [LGPL 2.1](https://github.com/moses-smt/mosesdecoder/blob/master/COPYING).
  Per questa raccolta adattata si applica l'opzione della sezione 3 di LGPL 2.1:
  GNU GPL versione 3, il cui testo è incluso nel file `LICENSE` del progetto.
- [spaCy, tokenizer_exceptions.py italiano](https://github.com/explosion/spaCy/blob/master/spacy/lang/it/tokenizer_exceptions.py):
  abbreviazioni ed eccezioni ortografiche sono separate. `cm` resta un simbolo di
  unità senza punto; apostrofi e sequenze di punti conservano le convenzioni del
  profilo Lingo. `C.so` è un'abbreviazione. Le divisioni `all'art.` e analoghe
  sono applicate nel profilo UD; il profilo standard conserva l'elisione intera.
  Non è importato l'insieme `BASE_EXCEPTIONS` di spaCy.

## spaCy — MIT

Copyright (C) 2016-2024 ExplosionAI GmbH, 2016 spaCy GmbH, 2015 Matthew Honnibal

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

## UD Italian-ISDT

LinGo was evaluated using the UD Italian-ISDT treebank, release r2.17:

- Repository: <https://github.com/UniversalDependencies/UD_Italian-ISDT>
- Release r2.17: <https://github.com/UniversalDependencies/UD_Italian-ISDT/tree/r2.17>
- License: <https://github.com/UniversalDependencies/UD_Italian-ISDT/blob/r2.17/LICENSE.txt>

The corpus files are not distributed with LinGo. The repository contains only
aggregate evaluation metrics, corpus checksums, and instructions for reproducing
the results. UD Italian-ISDT retains its own license and attribution requirements
independently of LinGo's source-code license.
