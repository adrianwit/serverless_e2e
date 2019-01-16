package cloud_function

import (
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/functions/metadata"
	"context"
	"github.com/adrianwit/serverless_e2e/cloud_function/fs"
	"github.com/viant/toolbox"
	"strings"
	"time"
)

type EventValue struct {
	CreateTime time.Time              `json:"createTime"`
	UpdateTime time.Time              `json:"updateTime"`
	Fields     map[string]interface{} `json:"fields"`
	Name       string                 `json:"name"`
}

type UpdateMask struct {
	FieldPaths []string `json:"fieldPaths"`
}

type FirestoreEvent struct {
	Value      *EventValue `json:"value"`
	OldValue   *EventValue `json:"oldValue"`
	UpdateMask *UpdateMask `json:"updateMask"`
}

/*

gcloud alpha functions deploy ModeratePost --entry-point ModeratePostFn  \
	--trigger-event providers/cloud.firestore/eventTypes/document.write \
	--trigger-resource  'projects/abstractdb-154a9/databases/(default)/documents/posts/{doc}' --runtime go111

*/

func ModeratePostFn(ctx context.Context, event FirestoreEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return err
	}
	if len(event.UpdateMask.FieldPaths) == 0 {
		return nil
	}
	fragments := strings.Split(meta.Resource.Name, "/")
	projectID := fragments[1]
	collection := fragments[5]
	documentID := fragments[6]
	return fs.RunTransaction(ctx, projectID, collection, documentID, func(ref *firestore.DocumentRef, transaction *firestore.Transaction) error {
		doc, err := transaction.Get(ref)
		if err != nil {
			return err
		}
		converter := toolbox.NewConverter("", "json")
		aPost := &Post{}
		if err = converter.AssignConverted(aPost, doc.Data()); err != nil {
			return err
		}
		if len(aPost.Comments) == 0 {
			return nil
		}
		if !aPost.Moderate() {
			return nil
		}
		var comments = make([]map[string]interface{}, 0)
		if err := converter.AssignConverted(&comments, aPost.Comments); err != nil {
			return err
		}
		return transaction.Set(ref, map[string]interface{}{
			"comments": comments,
		}, firestore.MergeAll)
	})
}

//Comment represents a post comment
type Comment struct {
	Text      string `json:"text"`
	Moderated bool   `json:"comments"`
}

//Post represents a post
type Post struct {
	Text     string     `json:"text"`
	Comments []*Comment `json:"comments"`
}

//Moderate moderates supplied post comments
func (p *Post) Moderate() bool {
	if len(p.Comments) == 0 {
		return false
	}
	moderated := false
	for i, comment := range p.Comments {
		if comment.Moderated {
			continue
		}
		p.Comments[i].Moderated = true
		p.Comments[i].Text = moderate(comment.Text)
		moderated = true
	}
	return moderated
}

func moderate(text string) string {
	fragments := strings.Split(text, " ")
	moderated := make([]string, len(fragments))
	sentenceBegin := true
	for i, fragment := range fragments {
		fragment = strings.ToLower(fragment)
		if sentenceBegin {
			moderated[i] = strings.Title(fragment)
			sentenceBegin = false
			continue
		}
		moderated[i] = fragment
		sentenceBegin = strings.Contains(fragment, ".")
	}
	return strings.Join(moderated, " ")
}
