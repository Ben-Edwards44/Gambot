import pygame
import src.graphics.graphics_const as graphics_const


def draw_pieces(window, pieces):
    #pieces is a list of piece.Piece objects
    for i in pieces:
        img = pygame.image.load(i.img_path)
        img = pygame.transform.scale(img, (i.img_width, i.img_height))

        window.blit(img, (i.draw_x, i.draw_y))


def draw_squares(window):
    #draw the background squares

    for i in range(8):
        for j in range(8):
            x = i * graphics_const.STEP_X
            y = j * graphics_const.STEP_Y

            if (i + j) % 2 == 0:
                colour = graphics_const.LIGHT_SQ_COLOUR
            else:
                colour = graphics_const.DARK_SQ_COLOUR

            pygame.draw.rect(window, colour, (x, y, graphics_const.STEP_X, graphics_const.STEP_Y))


def draw_board(window, pieces):
    draw_squares(window)
    draw_pieces(window, pieces)