echo "x=abcdefg" | run_mlr put '$y = substr($x, 0, 0)'
echo "x=abcdefg" | run_mlr put '$y = substr($x, 0, 7)'
echo "x=abcdefg" | run_mlr put '$y = substr($x, 1, 7)'
echo "x=abcdefg" | run_mlr put '$y = substr($x, 1, 6)'
echo "x=abcdefg" | run_mlr put '$y = substr($x, 2, 5)'
echo "x=abcdefg" | run_mlr put '$y = substr($x, 2, 3)'
echo "x=abcdefg" | run_mlr put '$y = substr($x, 3, 3)'
echo "x=abcdefg" | run_mlr put '$y = substr($x, 4, 3)'
echo "x=1234567" | run_mlr put '$y = substr($x, 2, 5)'
