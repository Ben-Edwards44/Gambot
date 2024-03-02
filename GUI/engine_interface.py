import subprocess


ENGINE_PATH = "C:\\Users\\Ben Edwards\\Documents\\Programming\\Python\\Projects\\chess-engine\\chess-engine\\chess-engine.exe"


class Engine:
    def __init__(self, debug):
        self.debug = debug

        self.process = subprocess.Popen(ENGINE_PATH, stdin=subprocess.PIPE, stdout=subprocess.PIPE)

        self.check_uci()

    def __del__(self):
        self.process.stdout.close()
        self.process.stdin.close()
        self.process.kill()

    def set_pos(self, move_list):
        #set the position to the startpos + the moves played in move_list
        if len(move_list) > 0:
            moves = " ".join(move_list)
            cmd = f"position startpos moves {moves}"
        else:
            cmd = "position startpos"

        self.send_cmd(cmd)

    def set_fen(self, fen):
        #set the position to the fen string given
        cmd = f'position fen "{fen}"'
        self.send_cmd(cmd)

    def get_perft_nodes(self, depth):
        #run perft and return number of nodes - assumes position has been set
        self.send_cmd(f"go perft {depth}")

        output = ""
        while len(output) <= 15 or output[:15] != "Nodes searched:":
            output = self.read_line()

        _, nodes = output.split(": ")

        return int(nodes)

    def get_legal_moves(self):
        #gets a list of the legal moves from the engine's divide perft - assumes position has been set
        #NOTE: this is not a UCI command, so a better solution should probably be made in future

        self.send_cmd("go perft 1")

        output = ""
        moves = []
        while len(output) <= 15 or output[:15] != "Nodes searched:":
            output = self.read_line()
            move, _ = output.split(":")

            if len(move) == 4 or len(move) == 5:  #promotions will have a length of 5 (a2a1r)
                moves.append(move)

        return moves

    def get_move(self, **kwargs):
        #assumes the position has been set
        self.send_args("go", kwargs)

        output = ""
        while len(output) < 8 or output[:8] != "bestmove":
            output = self.read_line()

        output = output.split(" ")
        move = output[1]  #output will look like: bestmove e2e4 ponder c7c5

        return move
    
    def new_game(self):
        self.send_cmd("ucinewgame")
        self.check_ready()

    def check_uci(self):
        self.send_cmd("uci")

        output = ""
        while output != "uciok":
            output = self.read_line()

            if output != "" and output[:2] != "id" and output[:6] != "option" and output != "uciok":
                raise Exception(f"UCI handshake failed")

    def check_ready(self):
        self.send_cmd("isready")

        output = self.read_line()

        if output != "readyok":
            raise Exception("Ready check failed")

    def read_line(self):
        #read a line of stdout

        if self.process.poll() is not None:
            raise Exception("UCI engine process finished")  #process has exited for some reason
        
        output = self.process.stdout.readline()
        text = output.decode().strip()

        if self.debug:
            print(f"Recieved: {text}")

        return text
    
    def send_args(self, cmd_name, kwarg_dict):
        #send args in the UCI required format (name1 value1 name2 value2...)
        args = ""
        for k, v in kwarg_dict.items():
            args += f" {k} {v}"

        self.send_cmd(f"{cmd_name}{args}")
        
    def send_cmd(self, cmd):
        if self.debug:
            print(f"Sending: {cmd}")

        b_cmd = bytes(f"{cmd}\n", encoding="utf-8")

        self.process.stdin.write(b_cmd)
        self.process.stdin.flush()