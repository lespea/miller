mlr --nidx put '$1 = sub($1, "ab(.)d(..)g",  "ab<<\1>>d<<\2>>g")' ./regtest/cases-pending-go-port/c-dsl-filter-pattern-action/0075.input
