import draw
import board
import player
import engine_interface
import graphics_const

import pygame


def init():
    draw.init()

    engine_process = engine_interface.Engine(debug=False)

    init_engine(engine_process)

    human_clock = board.Clock(graphics_const.PLAYER_TIME)
    engine_clock = board.Clock(graphics_const.ENGINE_TIME)

    board_obj = board.Board(graphics_const.START_FEN, engine_process, human_clock, engine_clock)

    human = player.HumanPlayer(board_obj, graphics_const.PLAYER_WHITE)
    engine = player.EnginePlayer(board_obj, not graphics_const.PLAYER_WHITE)

    return board_obj, human, engine


def init_engine(engine_process: engine_interface.Engine):
    engine_process.perform_handshake(True)
    engine_process.set_fen(graphics_const.START_FEN, [])


def get_start_colour():
    params = graphics_const.START_FEN.split(" ")
    white_to_move = params[1] == "w"

    return white_to_move


def get_move(board, human, engine, white_to_move):
    if white_to_move:
        clock = board.white_clock
    else:
        clock = board.black_clock

    clock.start_counting()

    if white_to_move == human.colour:
        move = human.get_move()
    else:
        move = engine.get_move()

    clock.stop_counting()

    return move


check_game_end = lambda board: len(board.engine.get_legal_moves()) == 0


def play_game():
    board_obj, human, engine = init()
    white_to_move = get_start_colour()

    draw.draw_board(board_obj)
    draw.draw_clocks(board_obj)

    playing = True
    while playing:
        move = get_move(board_obj, human, engine, white_to_move)

        board_obj.update(move)
        
        white_to_move = not white_to_move

        draw.draw_board(board_obj)
        draw.draw_clocks(board_obj)

        if check_game_end(board_obj):
            print("Game over")
            playing = False

        for e in pygame.event.get():
            if e.type == pygame.QUIT:
                human.kill_process()
                engine.kill_process()

                quit()


def main():
    play_again = True
    while play_again:
        play_game()

        play_again = input("Play again? (y/n) ") == "y"


if __name__ == "__main__":
    main()