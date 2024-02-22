import src.api.api as api
import src.graphics.main as graphics

from time import time
from random import randint
from subprocess import check_output


#NOTE: fens are from https://www.chessprogramming.org/Perft_Results
PERFT_FENS = ["rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
              "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - - -",
              "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - - -",
              "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
              "r2q1rk1/pP1p2pp/Q4n2/bbp1p3/Np6/1B3NBn/pPPP1PPP/R3K2R b KQ - 0 1",
              "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
              "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10"]

WHITE_PIECES = ["P", "N", "B", "R", "K", "Q"]
BLACK_PIECES = [i.lower() for i in WHITE_PIECES]

FILES = ["a", "b", "c", "d", "e", "f", "g", "h"]

FEN_FILEPATH = "data/equal_fens.txt"

SHOW_GRAPHICS = False


def parse_board(board_fen):
    ranks = board_fen.split("/")

    board = []
    for i in ranks:
        inx = 0
        rank = [0 for _ in range(8)]

        for x in i:
            if x in WHITE_PIECES:
                rank[inx] = WHITE_PIECES.index(x) + 1
                inx += 1
            elif x in BLACK_PIECES:
                rank[inx] = BLACK_PIECES.index(x) + 7
                inx += 1
            else:
                inx += int(x)

        board.append(rank)

    return board


def parse_ep(ep_target, white_to_move):
    #prev pawn double represents where the pawn is now; ep target is where the capturing pawn would go
    if ep_target == "-":
        return [-1, -1]
    
    x = int(ep_target[1])
    y = FILES.index(ep_target[0])

    if white_to_move:
        return [x + 1, y]
    else:
        return [x - 1, y]


def parse_fen(fen):
    board, move_colour, castle, ep, half_moves, full_moves = fen.split(" ")

    parsed_board = parse_board(board)

    white_to_move = move_colour == "w"

    white_k_castle = "K" in castle
    white_q_castle = "Q" in castle
    black_k_castle = "k" in castle
    black_q_castle = "q" in castle

    prev_pawn_double = parse_ep(ep, white_to_move)

    graphics.game_state.game_state_obj.load_from_dict({"board" : parsed_board,
                                                       "white_to_move" : white_to_move,
                                                       "white_king_castle" : white_k_castle,
                                                       "white_queen_castle" : white_q_castle,
                                                       "black_king_castle" : black_k_castle,
                                                       "black_queen_castle" : black_q_castle,
                                                       "prev_pawn_double" : prev_pawn_double})
    

def choose_fens(num):
    with open(FEN_FILEPATH, "r") as file:
        fens = file.read().split("\n")

    chosen = []
    for _ in range(num):
        inx = randint(0, len(fens) - 1)
        fen = fens.pop(inx)

        chosen.append(fen)

    return chosen


def run_engine(script_name):
    #run go engine
    output = check_output(f"{script_name}.exe").decode()
    
    return output


def check_win():
    #assumes that the graphics game state object has been updated

    api.send_data("check_win", graphics.game_state.game_state_obj)
    run_engine("chess-engine")

    terminal_state = api.load_check_win()

    return terminal_state


def play_game(white, black):
    #assumes that game state obj has been updated to the starting state of the game
    seen_pos = {}

    game_end = "not_terminal"
    while game_end == "not_terminal":
        #make move
        api.send_data("move_gen", graphics.game_state.game_state_obj)

        try:
            if graphics.game_state.game_state_obj.white_to_move:
                run_engine(white)
            else:
                run_engine(black)
        except Exception:
            game_end = "aborted"
            if graphics.game_state.game_state_obj.white_to_move:
                print(white)
            else:
                print(black)
                
            print(graphics.game_state.game_state_obj.board)

        #recieve the made moved
        new_state = api.load_game_state()
        graphics.game_state.game_state_obj.load_from_dict(new_state)

        #deal with draws by repetition
        hash_board = tuple(tuple(i) for i in graphics.game_state.game_state_obj.board)            
        if hash_board in seen_pos.keys():
            seen_pos[hash_board] += 1

            if seen_pos[hash_board] >= 3:
                #draw by repetition
                game_end = "draw"
                break
        else:
            seen_pos[hash_board] = 1

        if SHOW_GRAPHICS: 
            graphics.draw_board(graphics.game_state.game_state_obj.board)
            graphics.pygame.display.update()

        game_end = check_win()

    return game_end


def engine_game(engine1, engine2, num_games):
    #play 2 different versions of the engine

    engine1 = engine1.replace("/", "\\")  #for when the scripts are in a different directory. TODO: make it work for linux
    engine2 = engine2.replace("/", "\\")  #for when the scripts are in a different directory. TODO: make it work for linux

    graphics.game_state.init_game_state(None)

    if SHOW_GRAPHICS:
        graphics.init_graphics()

    fens = choose_fens(num_games)

    win1 = 0
    win2 = 0
    draw = 0
    aborted = 0
    for i, x in enumerate(fens):
        #update game state obj
        parse_fen(x)

        for j in range(2):
            #assign white/black
            if j == 0:
                white_player = engine1
                black_player = engine2
            else:
                white_player = engine2
                black_player = engine1

            game_end = play_game(white_player, black_player)
            
            if game_end == "draw":
                draw += 1
            elif game_end == "white_win":
                if white_player == engine1:
                    win1 += 1
                else:
                    win2 += 1
            elif game_end == "black_win":
                if black_player == engine1:
                    win1 += 1
                else:
                    win2 += 1
            else:
                aborted += 1

            print(f"Games Played: {i + 1}\n{engine1} wins: {win1}\n{engine2} wins: {win2}\nDraws: {draw}\nAborted : {aborted}\n")

    print(f"End result:\n{engine1} wins: {win1}\n{engine2} wins: {win2}\nDraws: {draw}\nAborted : {aborted}")


def get_time(engine):
    #assumes state pos obj has been updated with the position

    api.send_data("move_gen", graphics.game_state.game_state_obj)

    output = run_engine(engine)

    search_results = output.split("\n")
    times = [i.split(" ")[-1] for i in search_results]

    milisecs = []
    for time in times:
        if time == "":
            continue

        num = ""
        end = ""
        for i in time:
            if i in "0123456789.":
                num += i
            else:
                end += i

        num = float(num)

        if end == "ms":
            milisecs.append(num)
        elif end == "s":
            milisecs.append(num * 1000)
        elif end == "µs":
            milisecs.append(num / 1000)
        else:
            raise Exception(f"Invalid time {time}")
        
    return milisecs


def speed_test(engine1, engine2, num_games):
    engine1 = engine1.replace("/", "\\")  #for when the scripts are in a different directory. TODO: make it work for linux
    engine2 = engine2.replace("/", "\\")  #for when the scripts are in a different directory. TODO: make it work for linux

    graphics.game_state.init_game_state(None)

    fens = choose_fens(num_games)

    win1 = 0
    win2 = 0
    for i, x in enumerate(fens):
        #update game state obj
        parse_fen(x)

        time_depths1 = get_time(engine1)
        time_depths2 = get_time(engine2)

        if len(time_depths1) == len(time_depths2):
            #look at time taken for second to last depth (last full search)
            if len(time_depths1) == 1:
                time_depths1.append(0)
                time_depths2.append(0)
 
            if time_depths1[-2] < time_depths2[-2]:
                win1 += 1
            else:
                win2 += 1
        elif len(time_depths1) > len(time_depths2):
            win1 += 1
        else:
            win2 += 1

        print(f"Evalutated: {i + 1}\n{engine1} wins: {win1}\n{engine2} wins: {win2}")

    print(f"End Result:\n{engine1} wins: {win1}\n{engine2} wins: {win2}")


def perft(depth, test):
    #do performance test. NOTE: white to move is assumed
    start = time()
    graphics.game_state.init_game_state(None)

    for pos_num, fen in enumerate(PERFT_FENS):
        #update game state obj
        parse_fen(fen)

        api.send_data("perft", graphics.game_state.game_state_obj, perft_depth=depth, perft_test=test)
        
        output = run_engine("chess-engine")

        print(f"Position {pos_num + 1}:")
        print(output)

    end = time()

    print(f"Total time: {end - start :.3f}s")