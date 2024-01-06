class Piece:
    def __init__(self, name, img_path, x, y):
        self.name = name
        self.img_path = img_path

        self.act_x = x
        self.act_y = y

        self.draw_x = y
        self.draw_y = x

    act_pos = lambda self: (self.act_x, self.act_y)
    draw_pos = lambda self: (self.draw_x, self.draw_y)