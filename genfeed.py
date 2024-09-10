import uuid


for i in range(10):
    with open(f"test{i}.sh", "a+") as tfile:
        tfile.write(f"echo $(date) test{i}\n")
        for _ in range(100_000):
            tfile.write(
                'curl -s -X POST -H "Content-Type: application/json" -d \'{"' + str(uuid.uuid4()) + '":"' + str(uuid.uuid4()) + '"}\' http://localhost:53072/con/context1 > /dev/null \n'
            )
        tfile.write(f"echo $(date) test{i}\n")

    print(f'nohup time ./test{i}.sh &\n')
