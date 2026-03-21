(stdout "\e[33mcar\bt\e[0m =\t=\nfoo\rb\vfoo\a\n\"\v\r\n")
(run "(stdout \"bar\n\")")
(stderr "foo")
;(stdout (? true (stdout "is true") (stdout "is false")))
