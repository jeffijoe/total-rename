import findBoards from '../../Board-b/findBoards'

test('the board is the best', () => {
	findBoards().then(board => board.theBoardName)
})
