mlr head -n 4 then put -q -o json '@a[NR]=$a; @b[NR]=$b; emitp > "reg-test/cases-pending-go-port/c-dsl-redirects/0062.out.".$a.$b, (@a, @b), "NR"' reg-test/input/abixy
