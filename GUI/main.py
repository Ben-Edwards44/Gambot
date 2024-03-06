import draw
import utils
import input
import graphics_const
import engine_interface

import pygame


START_BOARD = [[10, 8, 9, 12, 11, 9, 8, 10], [7, 7, 7, 7, 7, 7, 7, 7], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [1, 1, 1, 1, 1, 1, 1, 1], [4, 2, 3, 6, 5, 3, 2, 4]]


def player_move(engine, board, move_list):
    engine.set_pos(move_list)
    legal_moves = engine.get_legal_moves()

    move = input.get_move(board, legal_moves)

    return move
        

def engine_move(engine, move_list):
    engine.set_pos(move_list)

    return engine.get_move()


def exit(engine):
    engine.kill_process()  #ensure the background engine process is killed

    quit()


def main():
    draw.init()

    move_list = []
    white_to_move = True

    engine = engine_interface.Engine()
    engine.new_game()

    board = START_BOARD

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