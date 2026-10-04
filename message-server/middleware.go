package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kaleb-white/letthemknow/message-server/utils"
)

// Requires slugs to be present and positive integers
func HasSlugs(slugs []string, h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wasError := false
		errs := []string{}

		for _, slug := range slugs {
			if s := r.PathValue(slug); s == "" {
				wasError = true					
				errs = append(errs, slug)
			} else if _, err := strconv.ParseUint(s, 10, 64); err != nil {
				wasError = true
				errs = append(errs, err.Error())
			}
		}

		if wasError {
			postfix := " not found in path"
		  s := strings.Builder{}
			numErrs := len(errs)

			// Single error
			if numErrs == 1 {
				fmt.Fprintf(&s, "%s%s", errs[0], postfix)
			} else {

				// Multi error
				for i, err := range errs {
					// Write until n - 1 
					if i < numErrs - 1 {
						fmt.Fprintf(&s, "%s, ", err)
					}
				}

				fmt.Fprintf(&s, "%s%s", errs[numErrs - 1], postfix)
			}

			utils.EncodeHttpError(s.String(), http.StatusBadRequest, w, r)
			return
		}

		h(w, r)
	})
}
