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

    step_x = graphics_const.SCREEN_WIDTH // 8
    step_y = graphics_const.SCREEN_HEIGHT // 8

    for i in range(8):
        for j in range(8):
            x = i * step_x
            y = j * step_y

            if (i + j) % 2 == 0:
                colour = graphics_const.LIGHT_SQ_COLOUR
            else:
                colour = graphics_const.DARK_SQ_COLOUR

            pygame.draw.rect(window, colour, (x, y, step_x, step_y))


def draw_board(window, pieces):
    draw_squares(window)
    draw_pieces(window, pieces)