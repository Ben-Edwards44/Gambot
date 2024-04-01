import engine_interface

from random import randint


FEN_FILEPATH = "../data/equal_fens.txt"
MOVE_TIME = 500


def choose_fens(num):
    with open(FEN_FILEPATH, "r") as file:
        fens = file.read().split("\n")

    chosen = []
    for _ in range(num):
        inx = randint(0, len(fens) - 1)
        fen = fens.pop(inx)

        chosen.append(fen)

    return chosen


def get_nodes(output):
    #gets the nodes searched

    args = output.split(" ")[1:]  #get rid of info...
    
    for i, x in enumerate(args):
        if x == "nodes":
            nodes = args[i + 1]
            break

    return int(nodes)


def engine_move(engine, fen):
    engine.set_fen(fen, [])

    engine.send_cmd(f"go movetime {MOVE_TIME}")

    output = ""
    nodes = []
    while len(output) < 8 or output[:8] != "bestmove":
        output = engine.read_line()

        if output[:8] != "bestmove":
            #we don't want to try and get the nodes for an output like "bestmove e2e4"
            node_num = get_nodes(output)
            nodes.append(node_num)

    return nodes


def get_best(engine1, engine2, fen):
    #return the nodes searched at the depth at which they differed

    nodes1 = engine_move(engine1, fen)[:-1]  #the last entry in unreliable
    nodes2 = engine_move(engine2, fen)[:-1]  #the last entry in unreliable

    common_depth = min(len(nodes1), len(nodes2))

    return sum(nodes1[:common_depth]), sum(nodes2[:common_depth])


def main(path1, path2, num):
    engine1 = engine_interface.Engine(path1)
    engine2 = engine_interface.Engine(path2)

    fens = choose_fens(num)

    nodes1 = 0
    nodes2 = 0
    for i, x in enumerate(fens):
        engine1.new_game()
        engine2.new_game()
        
        n1, n2 = get_best(engine1, engine2, x)

        nodes1 += n1
        nodes2 += n2

        print(f"Position: {i + 1}\n{path1} nodes searched: {nodes1}\n{path2} nodes searched: {nodes2}")

    print(f"Final Result\n{path1} nodes searched: {nodes1}\n{path2} nodes searched: {nodes2}")

    engine1.kill_process()
    engine2.kill_process()