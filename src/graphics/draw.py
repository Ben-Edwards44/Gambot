import pygame
import graphics_const


def init_draw():
    global window

    pygame.init()
    pygame.display.set_caption("Chess Engine")

    window = pygame.display.set_mode((graphics_const.SCREEN_WIDTH, graphics_const.SCREEN_HEIGHT))


def draw_pieces(pieces):
    #pieces is a list of piece.Piece objects

    step_x = graphics_const.SCREEN_WIDTH // 8
    step_y = graphics_const.SCREEN_HEIGHT // 8

    for i in pieces:
        img = pygame.image.load(i.img_path)
        img = pygame.transform.scale(img, (step_x, step_y))

        x = i.draw_x * step_x
        y = i.draw_y * step_y

        window.blit(img, (x, y))


def draw_squares():
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


def draw_board(pieces):
    window.fill((0, 0, 0))

    draw_squares()
    draw_pieces(pieces)

    pygame.display.update()