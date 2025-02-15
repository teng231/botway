package railway

import (
	"context"
	"fmt"
	"os"

	"github.com/abdfnx/botway/constants"
	"github.com/railwayapp/cli/entity"
	"github.com/railwayapp/cli/ui"
	"github.com/railwayapp/cli/uuid"
)

func (h *Handler) Delete(ctx context.Context, req *entity.CommandRequest) error {
	user, err := h.ctrl.GetUser(ctx)

	if err != nil {
		return err
	}

	if user.Has2FA {
		fmt.Print(constants.INFO_BACKGROUND.Render("INFO"))
		fmt.Println(constants.INFO_FOREGROUND.Render(" Your account has 2FA enabled, you must delete your project on the Railway Dashboard."))

		return nil
	}

	if len(req.Args) > 0 {
		// projectID provided as argument
		arg := req.Args[0]

		if uuid.IsValidUUID(arg) {
			project, err := h.ctrl.GetProject(ctx, arg)
			if err != nil {
				return err
			}

			return h.ctrl.DeleteProject(ctx, project.Id)
		}

		project, err := h.ctrl.GetProjectByName(ctx, arg)

		if err != nil {
			return err
		}

		return h.ctrl.DeleteProject(ctx, project.Id)
	}

	isLoggedIn, err := h.ctrl.IsLoggedIn(ctx)

	if err != nil {
		return err
	}

	if isLoggedIn {
		return h.deleteFromAccount(ctx, req)
	}

	return h.deleteFromID(ctx, req)
}

func (h *Handler) deleteFromAccount(ctx context.Context, req *entity.CommandRequest) error {
	projects, err := h.ctrl.GetProjects(ctx)
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		fmt.Print(constants.FAIL_BACKGROUND.Render("ERROR"))
		fmt.Println(constants.FAIL_FOREGROUND.Render(" No Projects found."))

		return nil
	}

	project, err := ui.PromptProjects(projects)

	if err != nil {
		return err
	}

	name := os.Args[2]

	if project.Name != name {
		fmt.Print(constants.FAIL_BACKGROUND.Render("ERROR"))
		fmt.Println(constants.FAIL_FOREGROUND.Render(" The project name typed doesn't match the selected project to be deleted."))

		return nil
	}

	fmt.Printf("🔥 Deleting project %s\n", ui.MagentaText(name))

	err = h.deleteById(ctx, project.Id)

	if err != nil {
		return err
	}

	fmt.Print(constants.SUCCESS_BACKGROUND.Render("SUCCESS"))
	fmt.Println(" " + name + " Deleted successfully")

	return nil
}

func (h *Handler) deleteFromID(ctx context.Context, req *entity.CommandRequest) error {
	projectID, err := ui.PromptText("Enter your project id")
	if err != nil {
		return err
	}

	project, err := h.ctrl.GetProject(ctx, projectID)

	if err != nil {
		return err
	}

	fmt.Printf("🔥 Deleting project %s\n", ui.MagentaText(project.Name))

	err = h.deleteById(ctx, project.Id)

	if err != nil {
		return err
	}

	fmt.Print(constants.SUCCESS_BACKGROUND.Render("SUCCESS"))
	fmt.Println(" " + project.Name + " Deleted successfully")

	return nil
}

func (h *Handler) deleteById(ctx context.Context, projectId string) error {
	err := h.ctrl.DeleteProject(ctx, projectId)

	if err != nil {
		return err
	}

	return nil
}
