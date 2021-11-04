build: .always
	make -C go build
check:
	make -C go check
install:
	make -C go install
	make -C man install
.always:
	true
