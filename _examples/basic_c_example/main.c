#include <stdio.h>
#include "libunifiedpush.h"
#include <unistd.h>
#include <string.h>

static void newMessage(char *instance, uint8_t *msg, int len, char *id)
{
	//this message can be deserialized here from something like json or whatever encoding you like

	// also note that the arguments to each of these input strings is freed after the function call
	// so if you need the data you should copy it somewhere

	puts("new message: ");
	write(stdout, msg, len);
	puts("\n");
}

static void newEndpoint(char *instance, char *endpoint)
{
	printf("new endpoint received: %s\n", endpoint);
}

static void unregistered(char *instance)
{
	printf("instance unregistered\n");
}

void upRegistration()
{
	struct UPRegister_return ret = UPRegister("");
	char *reason = ret.r1;
	int status = ret.r0;
	switch (status)
	{
	case 99:
		printf("up registration error happened: %s\n", reason);
		break;
	case UP_REGISTER_STATUS_NEW_ENDPOINT:
		printf("Will get new endpoint soon\n");
		break;
	case UP_REGISTER_STATUS_FAILED:
		printf("up registeratoin status failed: %s\n", reason);
		break;
	case UP_REGISTER_STATUS_REFUSED:
		printf("up registration refused %s\n", reason);
		break;
	}
	UPFreeString(reason);
}

void pickDistributors()
{
	struct UPGetDistributors_return ret = UPGetDistributors();
	char **fooarr = ret.r0;
	size_t length = ret.r1;

	char *selectedDist;

	if (length == 0)
	{
		printf("No Distributors found, exiting...\n");
		exit(1);
	}
	else if (length == 1)
	{
		selectedDist = fooarr[0];
		printf("Only one distributor, %s, avaiible picking that\n", selectedDist);
	}
	else
	{
		for (int i = 0; i < length; i++)
		{
			if (fooarr[i] == NULL)
			{
				break;
			}
			printf("%d. %s\n", i, fooarr[i]);
		}
		printf("pick a distributor:  ");
		fflush(stdout);
		unsigned int choice;
		scanf("%u", &choice);

		//theoretically do some bounds checking
		selectedDist = fooarr[choice];

		printf("distributor %s picked\n", selectedDist);
	}

	bool ok = UPSaveDistributor(selectedDist);
	if (!ok) {
		printf("unable to save distrib\n");
	}
	UPFreeStringArray(fooarr, length);
}

int main(int argc, char **argv)
{

	bool ok = UPInitializeAndCheck("cc.malhotra.karmanyaah.testapp.cgo", "C example app", *newMessage, *newEndpoint, *unregistered);
	if (ok) printf("successfully initialized notifications\n");

	if (strnlen(UPGetDistributor(), 1) == 0)
		pickDistributors();
	upRegistration();

	//do whatever your program does
	while (1)
		sleep(1);
}
