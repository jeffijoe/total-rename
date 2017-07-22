import findSpaces from '../../Space-b/findSpaces'

test('the space is the best', () => {
	findSpaces().then(space => space.theSpaceName)
})
