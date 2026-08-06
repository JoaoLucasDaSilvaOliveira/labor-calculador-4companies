package navigation

import usecase "labor-calculador-4companies/internal/application/usecase/company"

type WorkSpaceDeps struct {
	getCompanyByIdUC usecase.GetCompanyByIdUsecase
}

type WorkSpaceRouter struct {	
	Router //embeded
	deps WorkSpaceDeps
}

func NewWorkSpaceRouter(){} //TODO: create this constructor with it`s parent dependencies and it's own deps

func (wsr *WorkSpaceRouter)SelectCompany(companyID int) error {
	// wsr.deps.getCompanyByIdUC.Execute()
	// if err != wsr.Reset() {}
	//idea: create error page and a function of WorkSpaceRouter for this error handling
	return nil
}