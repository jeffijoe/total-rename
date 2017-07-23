import Promise from 'bluebird'
import apiHelper from '../../_helpers/apiHelper'
import BoardMemberRole from 'const/BoardMemberRole'
import { catchAsync } from '../../_helpers/catch'

describe('board access API', function () {
  let helper
  beforeEach(async function () {
    helper = await apiHelper({
      user: {
        firstName: 'Jeff',
        lastName: 'Hansen',
        username: `jeffijoe${Date.now()}`
      }
    })
  })

  describe('GET /boards/:id/members', function () {
    let member1, member2, member3, board
    beforeEach(async function () {
      const result = await Promise.all([
        apiHelper({ user: { firstName: 'Jon', lastName: 'West' } }),
        apiHelper({ user: { firstName: 'Amanda', lastName: 'Callesen' } }),
        apiHelper({ user: { firstName: 'Bjarke', lastName: 'Søgaard' } }),
        helper.createBoard()
      ])

      member1 = result[0]
      member2 = result[1]
      member3 = result[2]
      board = result[3]
      const arr = [member3, member1, member2] // ordering should not matter
      await Promise.all(
        arr.map(
          m => helper.createBoardAccess({
            boardId: board.id,
            userId: m.user.id
          })
            .then(
              () => m.acceptBoardInvite(board.id)
            )
        )
      )

      // Promote Jon to admin
      await helper.updateBoardAccess({
        boardId: board.id,
        userId: member1.user.id,
        role: BoardMemberRole.ADMIN
      })
    })

    it('returns the board members ordered', async function () {
      let [membersFromOwner, membersFrom1, membersFrom2] = await Promise.all([
        helper.getBoardMembers(board.id),
        member1.getBoardMembers(board.id),
        member2.getBoardMembers(board.id)
      ])

      expect(membersFromOwner).to.deep.equal(membersFrom1)
      expect(membersFromOwner).to.deep.equal(membersFrom2)
      expect(membersFromOwner.length).to.equal(4)
      const [access1, access2, access3, access4] = membersFromOwner
      expect(access1.user).to.exist
      expect(access1.user.email).to.not.exist
      expect(access1.user.password).to.not.exist
      expect(access1.user.username).to.exist

      expect(access1.user.firstName).to.equal(helper.user.firstName, 'first user should be the owner')
      expect(access2.user.firstName).to.equal(member1.user.firstName, 'second user should be the admin')
      expect(access3.user.firstName).to.equal(member2.user.firstName, 'third user should be based on name sorting')
      expect(access4.user.firstName).to.equal(member3.user.firstName, 'fourth user should be based on name sorting')
    })
  })

  describe('GET /boards/:boardId/members/:userId', function () {
    let board, other
    beforeEach(function () {
      return Promise.join(
        apiHelper(),
        helper.createBoard(),
        (o, s) => {
          other = o
          board = s
          return helper.addToBoard(board.id, other)
        }
      )
    })

    it('returns the board access', async function () {
      const ownerAccess1 = await helper.getBoardAccess(board.id, helper.user.id)
      const ownerAccess2 = await other.getBoardAccess(board.id, helper.user.id)
      expect(ownerAccess1).to.deep.equal(ownerAccess2)

      const memberAccess1 = await helper.getBoardAccess(board.id, other.user.id)
      const memberAccess2 = await other.getBoardAccess(board.id, other.user.id)

      expect(memberAccess1).to.deep.equal(memberAccess2)
      expect(memberAccess1.user).to.exist
      expect(memberAccess1.user.username).to.exist
      expect(memberAccess1.user.email).to.not.exist
      expect(memberAccess1.user.password).to.not.exist
    })

    describe('when not a member of the board', function () {
      it('throws a Forbidden error', async function () {
        const notAMember = await apiHelper()
        const { response } = await catchAsync(() => notAMember.getBoardAccess(board.id, helper.user.id))
        expect(response.status).to.equal(403)
      })
    })
  })
})
