import graphics_const

import pygame


class Piece:
    def __init__(self, window, x, y, piece_val):
        self.window = window

        self.x = x
        self.y = y

        self.img_width = graphics_const.STEP_X
        self.img_height = graphics_const.STEP_Y

        self.img_path = self.get_image_path(piece_val)
        self.draw_x, self.draw_y = self.get_draw_pos()

    def get_image_path(self, piece_val):
        if piece_val > 6:
            colour = "black"
            piece_name = graphics_const.PIECE_NAMES[piece_val - 7]
        else:
            colour = "white"
            piece_name = graphics_const.PIECE_NAMES[piece_val - 1]

        return f"images/{colour}/{piece_name}.png"

    def get_draw_pos(self): 
        #x and y are swapped because the array inxs are opposite to cartesian coords
        draw_x = self.y * graphics_const.STEP_Y + graphics_const.BOARD_TL[0]
        draw_y = self.x * graphics_const.STEP_X + graphics_const.BOARD_TL[1]

        return draw_x, draw_y

    def draw(self):
        img = pygame.image.load(self.img_path)
        img = pygame.transform.smoothscale(img, (self.img_width, self.img_height))

        self.window.blit(img, (self.draw_x, self.draw_y))


class DraggingPiece(Piece):
    def __init__(self, window, x, y, piece_val, mouse_x, mouse_y):
        self.mouse_x = mouse_x
        self.mouse_y = mouse_y

        super().__init__(window, x, y, piece_val)

    def get_draw_pos(self):
        #(OVERRIDE) we actually want to draw the piece to the mouse position
        offset_x = graphics_const.STEP_X // 2
        offset_y = graphics_const.STEP_Y // 2

        draw_x = self.mouse_x - offset_x
        draw_y = self.mouse_y - offset_y

        return draw_x, draw_y


def get_pieces(window, board_list, dragging_x=None, dragging_y=None):
    #returns a list of pieces (and maybe a dragging piece) from a board list
    pieces = []
    dragging_piece = None

    for i, x in enumerate(board_list):
        for j, k in enumerate(x):
            if k != 0:

                if i == dragging_x and j == dragging_y:
                    mouse_x, mouse_y = pygame.mouse.get_pos()
                    dragging_piece = DraggingPiece(window, i, j, k, mouse_x, mouse_y)
                else:
                    piece = Piece(window, i, j, k)
                    pieces.append(piece)

    return pieces, dragging_piece