resource "cloudflare_account_member" "member" {
  for_each = { for m in local.members : m.email => m }

  email_address = each.value.email
  role_ids      = [for r in each.value.roles : local.roles_map[r]]
}
