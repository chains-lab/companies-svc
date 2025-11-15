package middlewares

//func (s Service) CompanyMemberOrAdmin(
//	UserCtxKey interface{},
//	allowedCompanyRoles map[string]bool,
//	allowedAdminRoles map[string]bool,
//) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			ctx := r.Context()
//
//			user, ok := ctx.Value(UserCtxKey).(token.UserData)
//			if !ok {
//				ape.RenderErr(w,
//					problems.Unauthorized("Missing AuthorizationHeader header"),
//				)
//
//				return
//			}
//
//			if allowedAdminRoles[user.Role] {
//				next.ServeHTTP(w, r)
//				return
//			}
//
//			companyID, err := uuid.Parse(chi.URLParam(r, "company_id"))
//			if err != nil {
//				ape.RenderErr(w,
//					problems.BadRequest(validation.Errors{
//						"company_id": err,
//					})...,
//				)
//				return
//			}
//
//			if user.CompanyID == nil {
//				s.log.Error("user has no associated company", "user_id", user.ID)
//				ape.RenderErr(w, problems.Forbidden("User is not associated with any company"))
//				return
//			}
//
//			if companyID != *user.CompanyID {
//				s.log.Error("user company ID does not match", "user_id", user.ID, "user_company_id", user.CompanyID, "requested_company_id", companyID)
//				ape.RenderErr(w, problems.Forbidden("User does not belong to the requested company"))
//				return
//			}
//
//			if !(allowedCompanyRoles[user.Role]) {
//				s.log.Error("user role not allowed", "user_id", user.ID, "user_role", user.Role)
//				ape.RenderErr(w, problems.Forbidden("User does not have the required role"))
//				return
//			}
//
//			next.ServeHTTP(w, r)
//		})
//	}
//}
