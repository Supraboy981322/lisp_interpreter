(stdout "\e[33mcar\bt\e[0m =\t=\nfoo\rb\vfoo\a\n\"\v\r\n")
(run "(stdout \"bar\n\")")
(stderr "foo")
#|
 TODO: evaluate conditions
 (stdout (? true "is true" "is false"))
  #| nested comment |#
|#
