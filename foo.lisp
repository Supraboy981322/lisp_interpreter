;(stdout "\e[33mcar\bt\e[0m =\t=\nfoo\rb\vfoo\a\n\"\v\r\n")
;(run "(stdout \"bar\n\")")
;(stderr "\npoo" "\rf\n")
;(stderr "\nfarts" "\rp\n")
;(stderr "\nbar" "\bt\n" "\r")
(stdout (< 1 0))
#|
 TODO: evaluate conditions
 (stdout (? (< 1 0) "is true" "is false"))
  #| nested comment |#
|#
