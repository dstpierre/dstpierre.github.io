func render(w http.ResponseWriter, r *http.Request, view string, data any) {
  locale := r.Context().Value(middleware.ContextKeyLanguage).(string)
  lang := locale[:2]

  pdata := tpl.PageData{
    Lang:        lang,
    Locale:      locale,
    XSRFToken:   xsrftoken.Generate(XSRFToken, "", ""),
    Data:        data,
    CurrentUser: getCurrentUser(r),
  }

  if err := templ.Render(w, view, pdata); err != nil {
    fmt.Printf("error rendering %s -> %v", view, err)
    http.Redirect(w, r, "/error", http.StatusSeeOther)
  }
}
