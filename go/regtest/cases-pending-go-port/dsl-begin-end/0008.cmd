mlr --from regtest/input/s.dkvp put -q 'begin{} @sum[1+(NR%2)] += $x; end{dump}'
