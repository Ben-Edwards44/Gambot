import src.graphics.graphics_const as graphics_const


class Piece:
    def __init__(self, name, img_path, x, y):
        self.name = name
        self.img_path = img_path

        self.act_x = x
        self.act_y = y

        self.draw_x, self.draw_y = self.get_draw_pos(x, y)

    def get_draw_pos(self, x, y): 
        step_x = graphics_const.SCREEN_WIDTH // 8
        step_y = graphics_const.SCREEN_HEIGHT // 8

        #x and y are swapped because the array inxs are opposite to cartesian coords
        draw_x = y * step_y
        draw_y = x * step_x

        return draw_x, draw_y

    act_pos = lambda self: (self.act_x, self.act_y)
    draw_pos = lambda self: (self.draw_x, self.draw_y)