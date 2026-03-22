;(print[out] "\e[33mcar\bt\e[0m =\t=\nfoo\rb\vfoo\a\n\"\v\r\n")
;(run "(stdout \"bar\n\")")
;(print[err] "\npoo" "\rf\n")
;(print[out] "\nfarts" "\rp\n")
;(print[err] "\nbar" "\bt\n" "\r")

;(print[err] "foo\n")
(print[out] (? (< 1 0) "foo" "bar"))

#|
 TODO: evaluate conditions
 (print[out] (? (< 1 0) "is true" "is false"))
  #| nested comment |#
|#
