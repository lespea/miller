mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stderr, @*, "NR"' reg-test/input/abixy
