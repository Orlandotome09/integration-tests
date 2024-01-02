package posProcessor

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	BexsInternalBankCode = "144"
)

type posProcessor struct {
	partnerAdapter       interfaces.PartnerAdapter
	accountAdapter       interfaces.AccountAdapter
	addressAdapter       interfaces.AddressAdapter
	treeAdapterPublisher interfaces.TreeAdapterPublisher
	limitPublisher       interfaces.LimitPublisher
}

func NewPosProcessor(partnerAdapter interfaces.PartnerAdapter,
	accountAdapter interfaces.AccountAdapter,
	addressAdapter interfaces.AddressAdapter,
	treeAdapterPublisher interfaces.TreeAdapterPublisher,
	limitPublisher interfaces.LimitPublisher,
) interfaces.PosProcessor {
	return &posProcessor{
		partnerAdapter:       partnerAdapter,
		accountAdapter:       accountAdapter,
		addressAdapter:       addressAdapter,
		treeAdapterPublisher: treeAdapterPublisher,
		limitPublisher:       limitPublisher,
	}
}

func (ref *posProcessor) SendToTreeAdapter(profile entity.Profile, catalogConfig entity.ProductConfig) error {

	if profile.RoleType == values.RoleTypeCustomer && catalogConfig.TreeIntegration {
		partner, err := ref.partnerAdapter.GetActive(profile.PartnerID)
		if err != nil {
			return errors.WithStack(err)
		}

		profile.PartnerID = partner.Name

		accounts, err := ref.accountAdapter.FindByProfileID(profile.ProfileID.String())
		if err != nil {
			return errors.WithStack(err)
		}

		addresses, err := ref.addressAdapter.Search(profile.ProfileID.String())
		if err != nil {
			return errors.WithStack(err)
		}

		if err := ref.treeAdapterPublisher.Send(profile, accounts, addresses); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}

	logrus.Infof("[SendToTreeAdapter] Not sending profile to Temis Tree Adapter since it is not Enabled for this integration, or it is not a Customer Profile")

	return nil
}

func (ref *posProcessor) SendToLimit(profile entity.Profile, state entity.State, catalogConfig entity.ProductConfig) error {
	if catalogConfig.LimitIntegration {
		if err := ref.limitPublisher.Send(profile, state); err != nil {
			return errors.WithStack(err)
		}
	}

	logrus.Infof("[SendToLimit] Not sending profile to Temis Limit since it is not Enabled for this integration")

	return nil
}

func (ref *posProcessor) CreateInternalAccount(entityID uuid.UUID) error {
	logrus.Infof("[Profile Engine] About to create account for profile %v ", entityID.String())
	_, err := ref.accountAdapter.CreateInternal(entityID.String(), BexsInternalBankCode)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
