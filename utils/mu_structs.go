package utils

type MuLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MuLoginResponse struct {
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Context struct {
		SessionToken string `json:"session_token"`
	} `json:"context"`
}

type MuLogoutResponse struct {
	Status  string   `json:"status"`
	Reason  string   `json:"reason"`
	Context struct{} `json:"context"`
}

type MuSeriesInfoResponse struct {
	SeriesID   int    `json:"series_id"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	Associated []struct {
		Title string `json:"title"`
	} `json:"associated"`
	Description string `json:"description"`
	Image       struct {
		URL struct {
			Original string `json:"original"`
			Thumb    string `json:"thumb"`
		} `json:"url"`
		Height int `json:"height"`
		Width  int `json:"width"`
	} `json:"image"`
	Type           string  `json:"type"`
	Year           string  `json:"year"`
	BayesianRating float32 `json:"bayesian_rating"`
	RatingVotes    int     `json:"rating_votes"`
	Genres         []struct {
		Genre string `json:"genre"`
	} `json:"genres"`
	Categories []struct {
		SeriesID   int    `json:"series_id"`
		Category   string `json:"category"`
		Votes      int    `json:"votes"`
		VotesPlus  int    `json:"votes_plus"`
		VotesMinus int    `json:"votes_minus"`
		AddedBy    int    `json:"added_by"`
	} `json:"categories"`
	LatestChapter int    `json:"latest_chapter"`
	ForumID       int    `json:"forum_id"`
	Status        string `json:"status"`
	Licensed      bool   `json:"licensed"`
	Completed     bool   `json:"completed"`
	Anime         struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"anime"`
	RelatedSeries []struct {
		RelationID            int    `json:"relation_id"`
		RelationType          string `json:"relation_type"`
		RelatedSeriesID       int    `json:"related_series_id"`
		RelatedSeriesName     string `json:"related_series_name"`
		TriggeredByRelationID int    `json:"triggered_by_relation_id"`
	} `json:"related_series"`
	Authors []struct {
		Name     string `json:"name"`
		AuthorID int    `json:"author_id"`
		Type     string `json:"type"`
	} `json:"authors"`
	Publishers []struct {
		PublisherName string `json:"publisher_name"`
		PublisherID   int    `json:"publisher_id"`
		Type          string `json:"type"`
		Notes         string `json:"notes"`
	} `json:"publishers"`
	Publications []struct {
		PublicationName string `json:"publication_name"`
		PublisherName   string `json:"publisher_name"`
		PublisherID     int    `json:"publisher_id"`
	} `json:"publications"`
	Recommendations []struct {
		SeriesName string `json:"series_name"`
		SeriesID   int    `json:"series_id"`
		Weight     int    `json:"weight"`
	} `json:"recommendations"`
	CategoryRecommendations []struct {
		SeriesName string `json:"series_name"`
		SeriesID   int    `json:"series_id"`
		Weight     int    `json:"weight"`
	} `json:"category_recommendations"`
	Rank struct {
		Position struct {
			Week        int `json:"week"`
			Month       int `json:"month"`
			ThreeMonths int `json:"three_months"`
			SixMonths   int `json:"six_months"`
			Year        int `json:"year"`
		} `json:"position"`
		OldPosition struct {
			Week        int `json:"week"`
			Month       int `json:"month"`
			ThreeMonths int `json:"three_months"`
			SixMonths   int `json:"six_months"`
			Year        int `json:"year"`
		} `json:"old_position"`
		Lists struct {
			Reading    int `json:"reading"`
			Wish       int `json:"wish"`
			Complete   int `json:"complete"`
			Unfinished int `json:"unfinished"`
			Custom     int `json:"custom"`
		} `json:"lists"`
	} `json:"rank"`
	LastUpdated struct {
		Timestamp int    `json:"timestamp"`
		AsRFC3339 string `json:"as_rfc3339"`
		AsString  string `json:"as_string"`
	} `json:"last_updated"`
	Admin struct {
		AddedBy struct {
			UserID   int    `json:"user_id"`
			Username string `json:"username"`
			URL      string `json:"url"`
			Avatar   struct {
				ID     int    `json:"id"`
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"avatar"`
			TimeJoined struct {
				Timestamp int    `json:"timestamp"`
				AsRFC3339 string `json:"as_rfc3339"`
				AsString  string `json:"as_string"`
			} `json:"time_joined"`
			Signature     string `json:"signature"`
			ForumTitle    string `json:"forum_title"`
			FoldingAtHome bool   `json:"folding_at_home"`
			Profile       struct {
				Upgrade struct {
					Requested bool   `json:"requested"`
					Reason    string `json:"reason"`
				} `json:"upgrade"`
			} `json:"profile"`
			Stats struct {
				ForumPosts      int `json:"forum_posts"`
				AddedAuthors    int `json:"added_authors"`
				AddedGroups     int `json:"added_groups"`
				AddedPublishers int `json:"added_publishers"`
				AddedReleases   int `json:"added_releases"`
				AddedSeries     int `json:"added_series"`
			} `json:"stats"`
			UserGroup     string `json:"user_group"`
			UserGroupName string `json:"user_group_name"`
		} `json:"added_by"`
		Approved bool `json:"approved"`
	} `json:"admin"`
}

type MuSearchGroupsRequest struct {
	Search  string `json:"search"`
	PerPage int    `json:"perpage"`
}

type MuSearchGroupsResponse struct {
	TotalHits int                   `json:"total_hits"`
	Page      int                   `json:"page"`
	PerPage   int                   `json:"perpage"`
	Results   []MuSearchGroupsGroup `json:"results"`
}

type MuSearchGroupsGroup struct {
	Record  MuGroup `json:"record"`
	HitName string  `json:"hit_name"`
}

type MuSearchSeriesRequest struct {
	Search  string `json:"search"`
	PerPage int    `json:"per_page"`
}

type MuSearchSeriesResponse struct {
	TotalHits int `json:"total_hits"`
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	Results   []struct {
		Record struct {
			SeriesID    int    `json:"series_id"`
			Title       string `json:"title"`
			URL         string `json:"url"`
			Description string `json:"description"`
			Image       struct {
				URL struct {
					Original string `json:"original"`
					Thumb    string `json:"thumb"`
				} `json:"url"`
				Height int `json:"height"`
				Width  int `json:"width"`
			} `json:"image"`
			Type           string  `json:"type"`
			Year           string  `json:"year"`
			BayesianRating float64 `json:"bayesian_rating"`
			RatingVotes    int     `json:"rating_votes"`
			Genres         []struct {
				Genre string `json:"genre"`
			} `json:"genres"`
			LatestChapter int `json:"latest_chapter"`
			Rank          struct {
				Position struct {
					Week        int `json:"week"`
					Month       int `json:"month"`
					ThreeMonths int `json:"three_months"`
					SixMonths   int `json:"six_months"`
					Year        int `json:"year"`
				} `json:"position"`
				OldPosition struct {
					Week        int `json:"week"`
					Month       int `json:"month"`
					ThreeMonths int `json:"three_months"`
					SixMonths   int `json:"six_months"`
					Year        int `json:"year"`
				} `json:"old_position"`
				Lists struct {
					Reading    int `json:"reading"`
					Wish       int `json:"wish"`
					Complete   int `json:"complete"`
					Unfinished int `json:"unfinished"`
					Custom     int `json:"custom"`
				} `json:"lists"`
			} `json:"rank"`
			LastUpdated struct {
				Timestamp int64  `json:"timestamp"`
				AsRFC3339 string `json:"as_rfc3339"`
				AsString  string `json:"as_string"`
			} `json:"last_updated"`
			Admin struct {
				AddedBy struct {
					UserID   int    `json:"user_id"`
					Username string `json:"username"`
					URL      string `json:"url"`
					Avatar   struct {
						ID     int    `json:"id"`
						URL    string `json:"url"`
						Height int    `json:"height"`
						Width  int    `json:"width"`
					} `json:"avatar"`
					TimeJoined struct {
						Timestamp int64  `json:"timestamp"`
						AsRFC3339 string `json:"as_rfc3339"`
						AsString  string `json:"as_string"`
					} `json:"time_joined"`
					Signature     string `json:"signature"`
					ForumTitle    string `json:"forum_title"`
					FoldingAtHome bool   `json:"folding_at_home"`
					Profile       struct {
						Upgrade struct {
							Requested bool   `json:"requested"`
							Reason    string `json:"reason"`
						} `json:"upgrade"`
					} `json:"profile"`
					Stats struct {
						ForumPosts      int `json:"forum_posts"`
						AddedAuthors    int `json:"added_authors"`
						AddedGroups     int `json:"added_groups"`
						AddedPublishers int `json:"added_publishers"`
						AddedReleases   int `json:"added_releases"`
						AddedSeries     int `json:"added_series"`
					} `json:"stats"`
					UserGroup     string `json:"user_group"`
					UserGroupName string `json:"user_group_name"`
				} `json:"added_by"`
				Approved bool `json:"approved"`
			} `json:"admin"`
		} `json:"record"`
		HitTitle string `json:"hit_title"`
		Metadata struct {
			UserList struct {
				Series struct {
					ID    int    `json:"id"`
					Title string `json:"title"`
				} `json:"series"`
				ListID   int    `json:"list_id"`
				ListType string `json:"list_type"`
				ListIcon string `json:"list_icon"`
				Status   struct {
					Volume  int `json:"volume"`
					Chapter int `json:"chapter"`
				} `json:"status"`
				Priority  int `json:"priority"`
				TimeAdded struct {
					Timestamp int64  `json:"timestamp"`
					AsRFC3339 string `json:"as_rfc3339"`
					AsString  string `json:"as_string"`
				} `json:"time_added"`
				UserGenreHighlights []struct {
					Genre string `json:"genre"`
					Color string `json:"color"`
				} `json:"user_genre_highlights"`
			} `json:"user_list"`
		} `json:"metadata"`
	} `json:"results"`
}

type MuSeriesGroupsResponse struct {
	GroupList   []MuGroup `json:"group_list"`
	ReleaseList []struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Volume  string `json:"volume"`
		Chapter string `json:"chapter"`
		Groups  []struct {
			Name    string `json:"name"`
			GroupID int    `json:"group_id"`
			URL     string `json:"url"`
		} `json:"groups"`
		ReleaseDate string `json:"release_date"`
		TimeAdded   struct {
			Timestamp int    `json:"timestamp"`
			AsRFC3339 string `json:"as_rfc3339"`
			AsString  string `json:"as_string"`
		} `json:"time_added"`
	} `json:"release_list"`
}

type MuGroup struct {
	GroupID int    `json:"group_id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	Social  struct {
		Site     string `json:"site"`
		Facebook string `json:"facebook"`
		Twitter  string `json:"twitter"`
		IRC      struct {
			Channel string `json:"channel"`
			Server  string `json:"server"`
		} `json:"irc"`
		Forum   string `json:"forum"`
		Discord string `json:"discord"`
	} `json:"social"`
	Active  bool   `json:"active"`
	Notes   string `json:"notes"`
	AddedBy struct {
		UserID   int    `json:"user_id"`
		Username string `json:"username"`
		URL      string `json:"url"`
		Avatar   struct {
			ID     int    `json:"id"`
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"avatar"`
		TimeJoined struct {
			Timestamp int    `json:"timestamp"`
			AsRFC3339 string `json:"as_rfc3339"`
			AsString  string `json:"as_string"`
		} `json:"time_joined"`
		Signature     string `json:"signature"`
		ForumTitle    string `json:"forum_title"`
		FoldingAtHome bool   `json:"folding_at_home"`
		Profile       struct {
			Upgrade struct {
				Requested bool   `json:"requested"`
				Reason    string `json:"reason"`
			} `json:"upgrade"`
		} `json:"profile"`
		Stats struct {
			ForumPosts      int `json:"forum_posts"`
			AddedAuthors    int `json:"added_authors"`
			AddedGroups     int `json:"added_groups"`
			AddedPublishers int `json:"added_publishers"`
			AddedReleases   int `json:"added_releases"`
			AddedSeries     int `json:"added_series"`
		} `json:"stats"`
		UserGroup     string `json:"user_group"`
		UserGroupName string `json:"user_group_name"`
	} `json:"added_by"`
}

type MuGroupInfoResponse struct {
	MuGroup    `json:",inline"`
	Associated []struct {
		Name string `json:"name"`
	} `json:"associated"`
	Admin struct {
		Approved bool `json:"approved"`
		Hold     bool `json:"hold"`
	} `json:"admin"`
}
