from sys import argv


n = int(argv[1])

b = bin(n)[2:]
while len(b) < 64:
    b = f"0{b}"


#since lsb should be top left
b = b[::-1]


for x in range(8):
    line = ""
    for y in range(8):
        line += b[x * 8 + y]

    print(line)