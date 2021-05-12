#include <stdio.h>
#include "libunifiedpush.h"

int main(int argc, char** argv) {

	DBusInitialize("cc.malhotra.karmanyaah.test.cgo");
	printf("HI\n");

	char** fooarr = ListDistributors();
	printf("%lu\n", sizeof(fooarr));
	for(int i = 0; i < sizeof(fooarr) / 8; i++){
		printf("%s\n", fooarr[i]);
	}
}

