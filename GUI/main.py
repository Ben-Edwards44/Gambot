import draw
import utils
import input
import graphics_const
import engine_interface

import pygame


def player_move(engine, board, move_list):
    engine.set_fen(graphics_const.START_FEN, move_list)
    legal_moves = engine.get_legal_moves()

    move = input.get_move(board, legal_moves)

    return move
        

def engine_move(engine, move_list):
    engine.set_fen(graphics_const.START_FEN, move_list)

    return engine.get_move(movetime=graphics_const.ENGINE_MOVE_TIME)


def exit(engine):
    engine.kill_process()  #ensure the background engine process is killed

    quit()


def start_from_fen(fen):
    board = utils.fen_to_board(fen)
    white_move = fen.split(" ")[1] == "w"

    return board, white_move


def main():
    #TODO: add clocks etc.

    draw.init()

    engine = engine_interface.Engine(debug=True)
    engine.new_game()

    move_list = []
    board, white_to_move = start_from_fen(graphics_const.START_FEN)

    draw.draw_board(board)

    while True:
        if white_to_move == graphics_const.PLAYER_WHITE:
            move = player_move(engine, board, move_list)

            if move == "QUIT":
                exit(engine)
        else:
            move = engine_move(engine, move_list)

        board = utils.make_move(move, board)

        move_list.append(move)
        white_to_move = not white_to_move

        draw.draw_board(board)

        for e in pygame.event.get():
            if e.type == pygame.QUIT:
                exit(engine)


if __name__ == "__main__":
    main()