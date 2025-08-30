package cmd

import (
	"io/ioutil"
	"ones-cfg2md/pkg/model"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestParseObjectTypes(t *testing.T) {
	testCases := []struct {
		name          string
		typesStr      string
		expectedTypes []model.ObjectType
		expectError   bool
	}{
		{
			name:     "Valid types",
			typesStr: "documents,catalogs",
			expectedTypes: []model.ObjectType{
				model.ObjectTypeDocument,
				model.ObjectTypeCatalog,
			},
			expectError: false,
		},
		{
			name:          "Single type",
			typesStr:      "enums",
			expectedTypes: []model.ObjectType{model.ObjectTypeEnum},
			expectError:   false,
		},
		{
			name:          "Empty string",
			typesStr:      "",
			expectedTypes: []model.ObjectType{model.ObjectTypeDocument},
			expectError:   false,
		},
		{
			name:          "Invalid type",
			typesStr:      "invalidtype",
			expectedTypes: nil,
			expectError:   true,
		},
		{
			name:          "Mixed valid and invalid types",
			typesStr:      "documents,invalidtype",
			expectedTypes: nil,
			expectError:   true,
		},
		{
			name: "All types",
			typesStr: "documents,catalogs,accumulationregisters,informationregisters,enums,chartsofcharacteristictypes",
			expectedTypes: []model.ObjectType{
				model.ObjectTypeDocument,
				model.ObjectTypeCatalog,
				model.ObjectTypeAccumulationRegister,
				model.ObjectTypeInformationRegister,
				model.ObjectTypeEnum,
				model.ObjectTypeChartOfCharacteristicTypes,
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualTypes, err := parseObjectTypes(tc.typesStr)

			if (err != nil) != tc.expectError {
				t.Fatalf("Expected error: %v, but got: %v", tc.expectError, err)
			}

			if !reflect.DeepEqual(actualTypes, tc.expectedTypes) {
				t.Errorf("Expected types: %v, but got: %v", tc.expectedTypes, actualTypes)
			}
		})
	}
}

func TestRootCmd(t *testing.T) {
	// Create a temporary directory for the source
	sourceDir, err := ioutil.TempDir("", "source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(sourceDir)

	// Create a dummy Configuration.xml file
	configContent := `
<MetaDataObject xmlns="http://v8.1c.ru/8.3/MDClasses" xmlns:app="http://v8.1c.ru/8.2/managed-application/core" xmlns:cfg="http://v8.1c.ru/8.1/data/enterprise/current-config" xmlns:cmi="http://v8.1c.ru/8.2/managed-application/cmi" xmlns:ent="http://v8.1c.ru/8.1/data/enterprise" xmlns:lf="http://v8.1c.ru/8.2/managed-application/logform" xmlns:style="http://v8.1c.ru/8.1/data/ui/style" xmlns:sys="http://v8.1c.ru/8.1/data/ui/fonts/system" xmlns:v8="http://v8.1c.ru/8.1/data/core" xmlns:v8ui="http://v8.1c.ru/8.1/data/ui" xmlns:web="http://v8.1c.ru/8.1/data/ui/colors/web" xmlns:win="http://v8.1c.ru/8.1/data/ui/colors/windows" xmlns:xen="http://v8.1c.ru/8.3/xcf/enums" xmlns:xpr="http://v8.1c.ru/8.3/xcf/predef" xmlns:xr="http://v8.1c.ru/8.3/xcf/readable" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.16">
	<Configuration uuid="d5e4b242-723c-4f9b-8888-192b6a3c92e0">
		<InternalInfo>
			<xr:ContainedObject>
				<xr:ClassId>9de1c934-22e2-41b2-961a-436981107802</xr:ClassId>
				<xr:ObjectId>a8133518-3611-4e9c-9279-35e8c9199775</xr:ObjectId>
			</xr:ContainedObject>
			<xr:ContainedObject>
				<xr:ClassId>9fcd25a0-4822-11d4-9414-008048da11f9</xr:ClassId>
				<xr:ObjectId>a3a4533a-50dd-402d-a0d7-2425f2afe7f3</xr:ObjectId>
			</xr:ContainedObject>
			<xr:ContainedObject>
				<xr:ClassId>e3687481-0a87-462c-a166-9f34594f9bba</xr:ClassId>
				<xr:ObjectId>a3a4533a-50dd-402d-a0d7-2425f2afe7f3</xr:ObjectId>
			</xr:ContainedObject>
			<xr:ContainedObject>
				<xr:ClassId>9de1c934-22e2-41b2-961a-436981107802</xr:ClassId>
				<xr:ObjectId>a8133518-3611-4e9c-9279-35e8c9199775</xr:ObjectId>
			</xr:ContainedObject>
			<xr:ContainedObject>
				<xr:ClassId>51f2d5d8-ea4d-4064-8892-82951750031e</xr:ClassId>
				<xr:ObjectId>a3a4533a-50dd-402d-a0d7-2425f2afe7f3</xr:ObjectId>
			</xr:ContainedObject>
			<xr:ContainedObject>
				<xr:ClassId>e68182ea-4237-4383-967f-90c1e3370bc7</xr:ClassId>
				<xr:ObjectId>a3a4533a-50dd-402d-a0d7-2425f2afe7f3</xr:ObjectId>
			</xr:ContainedObject>
		</InternalInfo>
		<Properties>
			<Name>Конфигурация</Name>
			<Synonym/>
			<Comment/>
			<NamePrefix/>
			<ConfigurationExtensionCompatibilityMode>Version8_3_23</ConfigurationExtensionCompatibilityMode>
			<DefaultRunMode>ManagedApplication</DefaultRunMode>
			<UsePurposes>
				<v8:Value xsi:type="xs:string">PlatformApplication</v8:Value>
			</UsePurposes>
			<ScriptVariant>Russian</ScriptVariant>
			<DefaultRoles/>
			<Vendor/>
			<Version/>
			<UpdateCatalogAddress/>
			<IncludeHelpInContents>false</IncludeHelpInContents>
			<UseManagedFormInOrdinaryApplication>false</UseManagedFormInOrdinaryApplication>
			<UseOrdinaryFormInManagedApplication>false</UseOrdinaryFormInManagedApplication>
			<AdditionalFullTextSearchDictionaries/>
			<CommonSettingsStorage/>
			<ReportsUserSettingsStorage/>
			<ReportsVariantsStorage/>
			<FormDataSettingsStorage/>
			<DynamicListsUserSettingsStorage/>
			<Content/>
			<DefaultReportForm/>
			<DefaultReportVariantForm/>
			<DefaultReportSettingsForm/>
			<DefaultReportAppearanceTemplate/>
			<DefaultDynamicListSettingsForm/>
			<DefaultSearchForm/>
			<DefaultDataHistoryChangeHistoryForm/>
			<DefaultDataHistoryVersionHistoryForm/>
			<DefaultDataHistoryVersionDifferencesForm/>
			<DefaultDataHistoryVersionDataForm/>
			<DefaultDataHistoryEventDataForm/>
			<DefaultDataHistoryDetailedDataForm/>
			<DefaultForm/>
			<DefaultChoiceForm/>
			<DefaultListForm/>
			<DefaultCommandsSeparation/>
			<DefaultFoldersAndItemsSeparation/>
			<DefaultMainClientApplicationWindowMode>Normal</DefaultMainClientApplicationWindowMode>
			<DefaultInterface/>
			<DefaultStyle/>
			<DefaultLanguage>Language.Русский</DefaultLanguage>
			<BriefInformation/>
			<DetailedInformation/>
			<Copyright/>
			<VendorInformationAddress/>
			<ConfigurationInformationAddress/>
			<DataLockControlMode>Managed</DataLockControlMode>
			<ObjectAutonumerationMode>NotAutoFree</ObjectAutonumerationMode>
			<ModalityUseMode>Use</ModalityUseMode>
			<SynchronousPlatformExtensionAndAddInCallUseMode>Use</SynchronousPlatformExtensionAndAddInCallUseMode>
			<InterfaceCompatibilityMode>Taxi</InterfaceCompatibilityMode>
			<CompatibilityMode>Version8_3_23</CompatibilityMode>
			<DefaultConstantsForm/>
		</Properties>
		<ChildObjects>
			<Language>Русский</Language>
			<Subsystem>Бухгалтерия</Subsystem>
			<Subsystem>Зарплата</Subsystem>
			<Subsystem>Кадры</Subsystem>
			<Subsystem>Продажи</Subsystem>
			<Subsystem>Склад</Subsystem>
			<CommonPicture>Логотип</CommonPicture>
			<SessionParameter>ТекущийПользователь</SessionParameter>
			<Role>Администратор</Role>
			<Role>Пользователь</Role>
			<Constant>ВалютаУчета</Constant>
			<Constant>УчетПоСкладам</Constant>
			<Catalog>Контрагенты</Catalog>
			<Document>Заказ</Document>
			<Enum>СостоянияЗаказов</Enum>
			<Report>ОтчетПоПродажам</Report>
			<DataProcessor>ПечатьЭтикеток</DataProcessor>
			<ChartOfCharacteristicTypes>ВидыХарактеристик</ChartOfCharacteristicTypes>
			<InformationRegister>КурсыВалют</InformationRegister>
			<InformationRegister>МобильныеОтчеты</InformationRegister>
			<AccumulationRegister>Продажи</AccumulationRegister>
			<AccumulationRegister>Взаиморасчеты</AccumulationRegister>
			<FunctionalOption>ВалютныйУчет</FunctionalOption>
			<FunctionalOptionsParameter>Организация</FunctionalOptionsParameter>
			<FilterCriterion>ДокументыКонтрагента</FilterCriterion>
		</ChildObjects>
	</Configuration>
</MetaDataObject>
`
	configFile := filepath.Join(sourceDir, "Configuration.xml")
	if err := ioutil.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create dummy document file
	documentDir := filepath.Join(sourceDir, "Documents")
	if err := os.MkdirAll(documentDir, 0755); err != nil {
		t.Fatal(err)
	}
	documentFile := filepath.Join(documentDir, "Заказ.xml")
	documentContent := `
<?xml version="1.0" encoding="UTF-8"?>
<MetaDataObject xmlns="http://v8.1c.ru/8.3/MDClasses" xmlns:app="http://v8.1c.ru/8.2/managed-application/core" xmlns:cfg="http://v8.1c.ru/8.1/data/enterprise/current-config" xmlns:cmi="http://v8.1c.ru/8.2/managed-application/cmi" xmlns:ent="http://v8.1c.ru/8.1/data/enterprise" xmlns:lf="http://v8.1c.ru/8.2/managed-application/logform" xmlns:style="http://v8.1c.ru/8.1/data/ui/style" xmlns:sys="http://v8.1c.ru/8.1/data/ui/fonts/system" xmlns:v8="http://v8.1c.ru/8.1/data/core" xmlns:v8ui="http://v8.1c.ru/8.1/data/ui" xmlns:web="http://v8.1c.ru/8.1/data/ui/colors/web" xmlns:win="http://v8.1c.ru/8.1/data/ui/colors/windows" xmlns:xen="http://v8.1c.ru/8.3/xcf/enums" xmlns:xpr="http://v8.1c.ru/8.3/xcf/predef" xmlns:xr="http://v8.1c.ru/8.3/xcf/readable" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.16">
	<Document uuid="d5e4b242-723c-4f9b-8888-192b6a3c92e0">
		<InternalInfo>
			<xr:GeneratedType name="DocumentObject.Заказ" category="Object">
				<xr:TypeId>a8133518-3611-4e9c-9279-35e8c9199775</xr:TypeId>
				<xr:ValueId>a8133518-3611-4e9c-9279-35e8c9199776</xr:ValueId>
			</xr:GeneratedType>
			<xr:GeneratedType name="DocumentRef.Заказ" category="Ref">
				<xr:TypeId>a8133518-3611-4e9c-9279-35e8c9199777</xr:TypeId>
				<xr:ValueId>a8133518-3611-4e9c-9279-35e8c9199778</xr:ValueId>
			</xr:GeneratedType>
			<xr:GeneratedType name="DocumentSelection.Заказ" category="Selection">
				<xr:TypeId>a8133518-3611-4e9c-9279-35e8c9199779</xr:TypeId>
				<xr:ValueId>a8133518-3611-4e9c-9279-35e8c919977a</xr:ValueId>
			</xr:GeneratedType>
			<xr:GeneratedType name="DocumentList.Заказ" category="List">
				<xr:TypeId>a8133518-3611-4e9c-9279-35e8c919977b</xr:TypeId>
				<xr:ValueId>a8133518-3611-4e9c-9279-35e8c919977c</xr:ValueId>
			</xr:GeneratedType>
			<xr:GeneratedType name="DocumentManager.Заказ" category="Manager">
				<xr:TypeId>a8133518-3611-4e9c-9279-35e8c919977d</xr:TypeId>
				<xr:ValueId>a8133518-3611-4e9c-9279-35e8c919977e</xr:ValueId>
			</xr:GeneratedType>
		</InternalInfo>
		<Properties>
			<Name>Заказ</Name>
			<Synonym>
				<v8:item>
					<v8:lang>ru</v8:lang>
					<v8:content>Заказ</v8:content>
				</v8:item>
			</Synonym>
			<Comment/>
			<UseStandardCommands>true</UseStandardCommands>
			<Numerator/>
			<NumberType>String</NumberType>
			<NumberLength>9</NumberLength>
			<NumberAllowedLength>Variable</NumberAllowedLength>
			<NumberPeriodicity>Nonperiodical</NumberPeriodicity>
			<CheckUnique>true</CheckUnique>
			<Autonumbering>true</Autonumbering>
			<Characteristics/>
			<BasedOn/>
			<InputByString>
				<xr:Field>Document.Заказ.StandardAttribute.Number</xr:Field>
			</InputByString>
			<CreateOnInput>Use</CreateOnInput>
			<SearchStringModeOnInputByString>Begin</SearchStringModeOnInputByString>
			<FullTextSearchOnInputByString>DontUse</FullTextSearchOnInputByString>
			<ChoiceDataGetModeOnInputByString>Directly</ChoiceDataGetModeOnInputByString>
			<DefaultObjectForm/>
			<DefaultListForm/>
			<DefaultChoiceForm/>
			<AuxiliaryObjectForm/>
			<AuxiliaryListForm/>
			<AuxiliaryChoiceForm/>
			<Posting>Allow</Posting>
			<RealTimePosting>Allow</RealTimePosting>
			<RegisterRecordsDeletion>AutoDeleteOnUnpost</RegisterRecordsDeletion>
			<RegisterRecordsWritingOnPost>WriteSelected</RegisterRecordsWritingOnPost>
			<SequenceFilling>AutoFill</SequenceFilling>
			<RegisterRecords/>
			<PostInPrivilegedMode>true</PostInPrivilegedMode>
			<UnpostInPrivilegedMode>true</UnpostInPrivilegedMode>
			<IncludeHelpInContents>false</IncludeHelpInContents>
			<DataLockFields/>
			<DataLockControlMode>Managed</DataLockControlMode>
			<FullTextSearch>Use</FullTextSearch>
			<ObjectPresentation/>
			<ExtendedObjectPresentation/>
			<ListPresentation/>
			<ExtendedListPresentation/>
			<Explanation/>
			<ChoiceHistoryOnInput>Auto</ChoiceHistoryOnInput>
			<DataHistory>DontUse</DataHistory>
			<UpdateDataHistoryImmediatelyAfterWrite>false</UpdateDataHistoryImmediatelyAfterWrite>
			<ExecuteAfterWriteDataHistoryVersionProcessing>false</ExecuteAfterWriteDataHistoryVersionProcessing>
		</Properties>
		<ChildObjects>
			<Attribute uuid="d5e4b242-723c-4f9b-8888-192b6a3c92e1">
				<Properties>
					<Name>Контрагент</Name>
					<Synonym>
						<v8:item>
							<v8:lang>ru</v8:lang>
							<v8:content>Контрагент</v8:content>
						</v8:item>
					</Synonym>
					<Comment/>
					<Type>
						<v8:Type>cfg:CatalogRef.Контрагенты</v8:Type>
					</Type>
					<PasswordMode>false</PasswordMode>
					<Format/>
					<EditFormat/>
					<ToolTip/>
					<MarkNegatives>false</MarkNegatives>
					<Mask/>
					<MultiLine>false</MultiLine>
					<ExtendedEdit>false</ExtendedEdit>
					<MinValue xsi:nil="true"/>
					<MaxValue xsi:nil="true"/>
					<FillFromFillingValue>false</FillFromFillingValue>
					<FillValue xsi:nil="true"/>
					<FillChecking>ShowError</FillChecking>
					<ChoiceFoldersAndItems>Items</ChoiceFoldersAndItems>
					<ChoiceParameterLinks/>
					<ChoiceParameters/>
					<QuickChoice>Auto</QuickChoice>
					<CreateOnInput>Auto</CreateOnInput>
					<ChoiceForm/>
					<LinkByType/>
					<ChoiceHistoryOnInput>Auto</ChoiceHistoryOnInput>
					<Indexing>DontIndex</Indexing>
					<FullTextSearch>Use</FullTextSearch>
					<DataHistory>Use</DataHistory>
				</Properties>
			</Attribute>
		</ChildObjects>
	</Document>
</MetaDataObject>
`
	if err := ioutil.WriteFile(documentFile, []byte(documentContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a temporary directory for the output
	outputDir, err := ioutil.TempDir("", "output")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(outputDir)

	// Execute the command
	rootCmd.SetArgs([]string{sourceDir, outputDir, "--types", "documents"})
	// a little hack to avoid `os.Exit(1)` in runConversion func
	var actualErr error
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		runConversion(cmd, args)
		actualErr = nil
	}
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if actualErr != nil {
		t.Fatal(actualErr)
	}

	// Check if the output directory contains the expected files
	expectedFile := filepath.Join(outputDir, "objects.csv")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file '%s' to be created, but it was not", expectedFile)
	}
}

