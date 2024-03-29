START_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"


PLAYER_WHITE = True


LEGAL_FILTER = True


ENGINE_MOVE_TIME = 500


SCREEN_WIDTH = 800
SCREEN_HEIGHT = 600


BOARD_TL = (144, 44)
BOARD_X = 512  #should be multiple of 8
BOARD_Y = 512  #should be multiple of 8

STEP_X = BOARD_X // 8
STEP_Y = BOARD_Y // 8


BORDER_WIDTH = 18
BORDER_TL = tuple(i - BORDER_WIDTH for i in BOARD_TL)
BORDER_X = BOARD_X + BORDER_WIDTH * 2
BORDER_Y = BOARD_Y + BORDER_WIDTH * 2
BORDER_CORNER_RADIUS = 5


FONT_NAME = None
FONT_SIZE = 22
FONT_COLOUR = (255, 255, 255)
FONT_ANTIALIAS = True


LIGHT_SQ_COLOUR = (240, 217, 181)
DARK_SQ_COLOUR = (181, 136, 99)
LEGAL_MOVE_COLOUR = (128, 88, 55)
BORDER_COLOUR = (64, 62, 60)


PIECE_NAMES = [
    "pawn",
    "knight",
    "bishop",
    "rook",
    "king",
    "queen"
]


PIECE_VALUES = {
    "pawn" : 1,
    "knight" : 2,
    "bishop" : 3,
    "rook" : 4,
    "king" : 5,
    "queen" : 6
}


FILES = ["a", "b", "c", "d", "e", "f", "g", "h"]

WHITE_PIECES = ["P", "N", "B", "R", "K", "Q"]
BLACK_PIECES = [i.lower() for i in WHITE_PIECES]