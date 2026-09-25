# Porting area

The H2APK reference source can be ported into this project under `internal/`.

The current CLI deliberately keeps the first test small: it proves that VS Code can execute a real Android build rather than pretending that a build succeeded.

Do not add fake APK generation or fake compiler output.
