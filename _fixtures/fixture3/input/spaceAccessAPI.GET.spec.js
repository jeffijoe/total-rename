import Promise from 'bluebird'
import apiHelper from '../../_helpers/apiHelper'
import SpaceMemberRole from 'const/SpaceMemberRole'
import { catchAsync } from '../../_helpers/catch'

describe('space access API', function () {
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

  describe('GET /spaces/:id/members', function () {
    let member1, member2, member3, space
    beforeEach(async function () {
      const result = await Promise.all([
        apiHelper({ user: { firstName: 'Jon', lastName: 'West' } }),
        apiHelper({ user: { firstName: 'Amanda', lastName: 'Callesen' } }),
        apiHelper({ user: { firstName: 'Bjarke', lastName: 'Søgaard' } }),
        helper.createSpace()
      ])

      member1 = result[0]
      member2 = result[1]
      member3 = result[2]
      space = result[3]
      const arr = [member3, member1, member2] // ordering should not matter
      await Promise.all(
        arr.map(
          m => helper.createSpaceAccess({
            spaceId: space.id,
            userId: m.user.id
          })
            .then(
              () => m.acceptSpaceInvite(space.id)
            )
        )
      )

      // Promote Jon to admin
      await helper.updateSpaceAccess({
        spaceId: space.id,
        userId: member1.user.id,
        role: SpaceMemberRole.ADMIN
      })
    })

    it('returns the space members ordered', async function () {
      let [membersFromOwner, membersFrom1, membersFrom2] = await Promise.all([
        helper.getSpaceMembers(space.id),
        member1.getSpaceMembers(space.id),
        member2.getSpaceMembers(space.id)
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

  describe('GET /spaces/:spaceId/members/:userId', function () {
    let space, other
    beforeEach(function () {
      return Promise.join(
        apiHelper(),
        helper.createSpace(),
        (o, s) => {
          other = o
          space = s
          return helper.addToSpace(space.id, other)
        }
      )
    })

    it('returns the space access', async function () {
      const ownerAccess1 = await helper.getSpaceAccess(space.id, helper.user.id)
      const ownerAccess2 = await other.getSpaceAccess(space.id, helper.user.id)
      expect(ownerAccess1).to.deep.equal(ownerAccess2)

      const memberAccess1 = await helper.getSpaceAccess(space.id, other.user.id)
      const memberAccess2 = await other.getSpaceAccess(space.id, other.user.id)

      expect(memberAccess1).to.deep.equal(memberAccess2)
      expect(memberAccess1.user).to.exist
      expect(memberAccess1.user.username).to.exist
      expect(memberAccess1.user.email).to.not.exist
      expect(memberAccess1.user.password).to.not.exist
    })

    describe('when not a member of the space', function () {
      it('throws a Forbidden error', async function () {
        const notAMember = await apiHelper()
        const { response } = await catchAsync(() => notAMember.getSpaceAccess(space.id, helper.user.id))
        expect(response.status).to.equal(403)
      })
    })
  })
})
