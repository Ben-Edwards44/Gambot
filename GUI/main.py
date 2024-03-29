import draw
import board
import player
import engine_interface
import graphics_const

import pygame


def init():
    draw.init()

    engine_process = engine_interface.Engine()

    init_engine(engine_process)

    board_obj = board.Board(graphics_const.START_FEN, engine_process)

    human = player.HumanPlayer(board_obj, graphics_const.PLAYER_WHITE, 0, 0)
    engine = player.EnginePlayer(board_obj, not graphics_const.PLAYER_WHITE, 0, 0)

    return board_obj, human, engine


def init_engine(engine_process: engine_interface.Engine):
    engine_process.check_uci()
    engine_process.new_game()
    engine_process.check_ready()
    engine_process.set_fen(graphics_const.START_FEN, [])


def get_start_colour():
    params = graphics_const.START_FEN.split(" ")
    white_to_move = params[1] == "w"

    return white_to_move


def get_move(human, engine, white_to_move):
    if white_to_move == human.colour:
        move = human.get_move()
    else:
        move = engine.get_move()

    return move


def main():
    #TODO: add clocks etc.

    board_obj, human, engine = init()
    white_to_move = get_start_colour()

    draw.draw_board(board_obj)

    while True:
        move = get_move(human, engine, white_to_move)

        board_obj.update(move)
        
        white_to_move = not white_to_move

        draw.draw_board(board_obj)

        for e in pygame.event.get():
            if e.type == pygame.QUIT:
                human.kill_process()
                engine.kill_process()

                quit()


if __name__ == "__main__":
    main()