import utils
import graphics_const

import pygame
from os import listdir


class Piece:
    def __init__(self, x, y, piece_val):
        self.act_x = x
        self.act_y = y

        self.img_path = images[piece_val]

        self.img_width = graphics_const.STEP_X
        self.img_height = graphics_const.STEP_Y

        self.draw_x, self.draw_y = self.get_draw_pos(x, y)

    def get_draw_pos(self, x, y): 
        #x and y are swapped because the array inxs are opposite to cartesian coords
        draw_x = y * graphics_const.STEP_Y
        draw_y = x * graphics_const.STEP_X

        return draw_x, draw_y
    
    def overwrite_draw_pos(self, mouse_x, mouse_y):
        #ensure center of image goes to mouse pos
        offset_x = graphics_const.STEP_X // 2
        offset_y = graphics_const.STEP_Y // 2

        self.draw_x = mouse_x - offset_x
        self.draw_y = mouse_y - offset_y


def init():
    global images
    global window

    images = get_images()

    pygame.init()
    pygame.display.set_caption("Chess Engine")

    window = pygame.display.set_mode((graphics_const.SCREEN_WIDTH, graphics_const.SCREEN_HEIGHT))


def get_images():
    #Get the paths to the images for each piece value

    images = {}
    for colour in ["white", "black"]:
        path = f"images/{colour}"
        img_names = listdir(path)

        for img_name in img_names:
            img_path = f"{path}/{img_name}"

            piece_name = img_name.split(".")[0]  #ignore the .pgn
            piece_value = graphics_const.PIECE_VALUES[piece_name]

            if colour == "black":
                piece_value += 6

            images[piece_value] = img_path

    return images


def draw_piece(piece):
    #piece is a Piece object

    img = pygame.image.load(piece.img_path)
    img = pygame.transform.scale(img, (piece.img_width, piece.img_height))

    window.blit(img, (piece.draw_x, piece.draw_y))


def draw_square(x, y, colour):
    draw_x = x * graphics_const.STEP_X
    draw_y = y * graphics_const.STEP_Y

    pygame.draw.rect(window, colour, (draw_x, draw_y, graphics_const.STEP_X, graphics_const.STEP_Y))


def draw_background():
    #draw the background squares

    for i in range(8):
        for j in range(8):
            if (i + j) % 2 == 0:
                colour = graphics_const.LIGHT_SQ_COLOUR
            else:
                colour = graphics_const.DARK_SQ_COLOUR

            draw_square(i, j, colour)


def draw_legal_moves(legal_moves, x, y):
    #colour the squares that the player could move to

    start_square = utils.square_to_str(x, y)

    for i in legal_moves:
        if i[:2] == start_square:
            end_y, end_x = utils.str_to_square(i[2:])  #swap x and y because the array inxs are opposite to cartesian coords

            draw_square(end_x, end_y, graphics_const.LEGAL_MOVE_COLOUR)


def draw_dragging_piece(board, selected_x, selected_y, legal_moves):
    #draw the board normally, apart from the dragging piece and the legal moves

    draw_background()
    draw_legal_moves(legal_moves, selected_x, selected_y)

    for i, x in enumerate(board):
        for j, k in enumerate(x):
            if k == 0:
                continue

            piece = Piece(i, j, k)

            if i == selected_x and j == selected_y:
                m_x, m_y = pygame.mouse.get_pos()
                piece.overwrite_draw_pos(m_x, m_y)

                dragging_piece = piece
            else:
                draw_piece(piece)

    draw_piece(dragging_piece)  #we want to draw dragging piece last so it is in front of other pieces

    pygame.display.update()


def draw_board(board):
    draw_background()

    for i, x in enumerate(board):
        for j, k in enumerate(x):
            if k == 0:
                continue
            
            piece = Piece(i, j, k)

            draw_piece(piece)

    pygame.display.update()