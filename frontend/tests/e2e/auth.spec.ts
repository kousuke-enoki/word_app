import { test, type UserRole } from './fixtures/auth'

test('signs in from /sign_in and lands on mypage', async ({
  goToMyPageFromSignIn,
}) => {
  await goToMyPageFromSignIn('admin')
})

test.describe('role restricted routes', () => {
  const roles: UserRole[] = ['general', 'admin', 'root']

  roles.forEach((role) => {
    test(`direct navigation respects ${role} permissions`, async ({
      expectRoleProtectedNavigation,
    }) => {
      await expectRoleProtectedNavigation(role)
    })
  })
})
