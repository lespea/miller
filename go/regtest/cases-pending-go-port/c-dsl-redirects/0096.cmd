mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "regtest/cases-pending-go-port/c-dsl-redirects/0096.out.".$a.$b, (@a, @b)' regtest/input/abixy
