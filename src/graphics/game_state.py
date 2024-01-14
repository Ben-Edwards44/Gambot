class gameState:
    def __init__(self, start_board):
        #for now, initialise defalut values

        self.board = start_board

        self.white_to_move = True
        
        self.white_can_castle = True
        self.black_can_castle = True

        #will store [x, y] pos of double pawn advance
        self.prev_pawn_double = [-1, -1]

    def get_dict(self):
        attr_dict = {
            "board" : self.board,
            "white_to_move" : self.white_to_move,
            "white_can_castle" : self.white_can_castle,
            "black_can_castle" : self.black_can_castle,
            "prev_pawn_double" : self.prev_pawn_double
        }

        return attr_dict
    
    def load_from_dict(self, attr_dict):
        #requires api data to be parsed (not all strings)

        self.board = attr_dict["board"]
        self.white_to_move = attr_dict["white_to_move"]
        self.white_can_castle = attr_dict["white_can_castle"]
        self.black_can_castle = attr_dict["black_can_castle"]
        self.prev_pawn_double = attr_dict["prev_pawn_double"]


def init_game_state(start_board):
    global game_state_obj

    game_state_obj = gameState(start_board)