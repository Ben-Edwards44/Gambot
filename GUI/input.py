import draw
import utils
import graphics_const

import pygame


def get_cell_inx():
    x, y = pygame.mouse.get_pos()

    space_x = graphics_const.SCREEN_WIDTH // 8
    space_y = graphics_const.SCREEN_HEIGHT // 8

    #x and y swap because the array inx is different to cartesian coords
    cell_x = y // space_y
    cell_y = x // space_x

    return cell_x, cell_y


def drag_piece(board, selected_x, selected_y, legal_moves):
    if selected_x == -1 and selected_y == -1:
        #player has just selected a piece
        x, y = get_cell_inx()

        #TODO: also check if it is correct colour
        if board[x][y] != 0:
            selected_x = x
            selected_y = y
    else:
        #player is moving selected piece
        draw.draw_dragging_piece(board, selected_x, selected_y, legal_moves)

    return selected_x, selected_y


def get_move(board, legal_moves):
    #get the player's move

    selected_x = -1
    selected_y = -1

    while True:        
        if pygame.mouse.get_pressed()[0]:
            selected_x, selected_y = drag_piece(board, selected_x, selected_y, legal_moves)
        else:
            if selected_x != -1 and selected_y != -1:
                end_x, end_y = get_cell_inx()
                move = utils.move_to_str(selected_x, selected_y, end_x, end_y)

                if move in legal_moves:
                    return move

            selected_x = -1
            selected_y = -1

            draw.draw_board(board)  #we need to draw the board after the player has let go of dragging piece

        for e in pygame.event.get():
            if e.type == pygame.QUIT:
                return "QUIT"