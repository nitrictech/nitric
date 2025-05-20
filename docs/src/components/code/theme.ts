import { Theme } from '@code-hike/lighter'

const theme: Theme = {
  name: 'Catppuccin',
  type: 'from-css',
  colors: {
    focusBorder: 'var(--code-mauve)',
    foreground: 'var(--code-text)',
    disabledForeground: 'var(--code-subtext-0)',
    'widget.shadow': 'var(--code-mantle)80',
    'selection.background': 'var(--code-mauve)66',
    descriptionForeground: 'var(--code-text)',
    errorForeground: 'var(--code-red)',
    'icon.foreground': 'var(--code-mauve)',
    'sash.hoverBorder': 'var(--code-mauve)',
    'window.activeBorder': 'var(--code-black)',
    'window.inactiveBorder': 'var(--code-black)',
    'textBlockQuote.background': 'var(--code-mantle)',
    'textBlockQuote.border': 'var(--code-crust)',
    'textCodeBlock.background': 'var(--code-base)',
    'textLink.activeForeground': 'var(--code-sky)',
    'textLink.foreground': 'var(--code-blue)',
    'textPreformat.foreground': 'var(--code-text)',
    'textSeparator.foreground': 'var(--code-mauve)',
    'activityBar.background': 'var(--code-crust)',
    'activityBar.foreground': 'var(--code-mauve)',
    'activityBar.dropBorder': 'var(--code-mauve)33',
    'activityBar.inactiveForeground': 'var(--code-overlay-0)',
    'activityBar.border': 'var(--code-black)',
    'activityBarBadge.background': 'var(--code-mauve)',
    'activityBarBadge.foreground': 'var(--code-crust)',
    'activityBar.activeBorder': 'var(--code-black)',
    'activityBar.activeBackground': 'var(--code-black)',
    'activityBar.activeFocusBorder': 'var(--code-black)',
    'activityBarTop.foreground': 'var(--code-mauve)',
    'activityBarTop.activeBorder': 'var(--code-black)',
    'activityBarTop.inactiveForeground': 'var(--code-overlay-0)',
    'activityBarTop.dropBorder': 'var(--code-mauve)33',
    'badge.background': 'var(--code-surface-1)',
    'badge.foreground': 'var(--code-text)',
    'breadcrumb.activeSelectionForeground': 'var(--code-mauve)',
    'breadcrumb.background': 'var(--code-base)',
    'breadcrumb.focusForeground': 'var(--code-mauve)',
    'breadcrumb.foreground': 'var(--code-text)cc',
    'breadcrumbPicker.background': 'var(--code-mantle)',
    'button.background': 'var(--code-mauve)',
    'button.foreground': 'var(--code-crust)',
    'button.border': 'var(--code-black)',
    'button.separator': 'var(--code-black)',
    'button.hoverBackground': '#dac1f9',
    'button.secondaryForeground': 'var(--code-text)',
    'button.secondaryBackground': 'var(--code-surface-2)',
    'button.secondaryHoverBackground': '#6a708c',
    'checkbox.background': 'var(--code-surface-1)',
    'checkbox.border': 'var(--code-black)',
    'checkbox.foreground': 'var(--code-mauve)',
    'dropdown.background': 'var(--code-mantle)',
    'dropdown.listBackground': 'var(--code-surface-2)',
    'dropdown.border': 'var(--code-mauve)',
    'dropdown.foreground': 'var(--code-text)',
    'debugToolBar.background': 'var(--code-crust)',
    'debugToolBar.border': 'var(--code-black)',
    'debugExceptionWidget.background': 'var(--code-crust)',
    'debugExceptionWidget.border': 'var(--code-mauve)',
    'debugTokenExpression.number': 'var(--code-peach)',
    'debugTokenExpression.boolean': 'var(--code-mauve)',
    'debugTokenExpression.string': 'var(--code-green)',
    'debugTokenExpression.error': 'var(--code-red)',
    'debugIcon.breakpointForeground': 'var(--code-red)',
    'debugIcon.breakpointDisabledForeground': 'var(--code-red)99',
    'debugIcon.breakpointUnverifiedForeground': '#a47487',
    'debugIcon.breakpointCurrentStackframeForeground': 'var(--code-surface-2)',
    'debugIcon.breakpointStackframeForeground': 'var(--code-surface-2)',
    'debugIcon.startForeground': 'var(--code-green)',
    'debugIcon.pauseForeground': 'var(--code-blue)',
    'debugIcon.stopForeground': 'var(--code-red)',
    'debugIcon.disconnectForeground': 'var(--code-surface-2)',
    'debugIcon.restartForeground': 'var(--code-teal)',
    'debugIcon.stepOverForeground': 'var(--code-mauve)',
    'debugIcon.stepIntoForeground': 'var(--code-text)',
    'debugIcon.stepOutForeground': 'var(--code-text)',
    'debugIcon.continueForeground': 'var(--code-green)',
    'debugIcon.stepBackForeground': 'var(--code-surface-2)',
    'debugConsole.infoForeground': 'var(--code-blue)',
    'debugConsole.warningForeground': 'var(--code-peach)',
    'debugConsole.errorForeground': 'var(--code-red)',
    'debugConsole.sourceForeground': 'var(--code-rosewater)',
    'debugConsoleInputIcon.foreground': 'var(--code-text)',
    'diffEditor.border': 'var(--code-surface-2)',
    'diffEditor.insertedTextBackground': 'var(--code-green)1a',
    'diffEditor.removedTextBackground': 'var(--code-red)1a',
    'diffEditor.insertedLineBackground': 'var(--code-green)26',
    'diffEditor.removedLineBackground': 'var(--code-red)26',
    'diffEditor.diagonalFill': 'var(--code-surface-2)99',
    'diffEditorOverview.insertedForeground': 'var(--code-green)cc',
    'diffEditorOverview.removedForeground': 'var(--code-red)cc',
    'editor.background': '#272140ff',
    'editor.findMatchBackground': '#604456',
    'editor.findMatchBorder': 'var(--code-red)33',
    'editor.findMatchHighlightBackground': '#455c6d',
    'editor.findMatchHighlightBorder': 'var(--code-sky)33',
    'editor.findRangeHighlightBackground': '#455c6d',
    'editor.findRangeHighlightBorder': 'var(--code-sky)33',
    'editor.foldBackground': 'var(--code-sky)40',
    'editor.foreground': 'var(--code-text)',
    'editor.hoverHighlightBackground': 'var(--code-sky)40',
    'editor.lineHighlightBackground': 'var(--code-text)12',
    'editor.lineHighlightBorder': 'var(--code-black)',
    'editor.rangeHighlightBackground': '#00000040',
    'editor.rangeHighlightBorder': 'var(--code-black)',
    'editor.selectionBackground': 'var(--code-overlay-2)40',
    'editor.selectionHighlightBackground': 'var(--code-overlay-2)33',
    'editor.selectionHighlightBorder': 'var(--code-overlay-2)33',
    'editor.wordHighlightBackground': 'var(--code-overlay-2)33',
    'editorBracketMatch.background': 'var(--code-overlay-2)1a',
    'editorBracketMatch.border': 'var(--code-overlay-2)',
    'editorCodeLens.foreground': 'var(--code-overlay-1)',
    'editorCursor.background': 'var(--code-base)',
    'editorCursor.foreground': 'var(--code-rosewater)',
    'editorGroup.border': 'var(--code-surface-2)',
    'editorGroup.dropBackground': 'var(--code-mauve)33',
    'editorGroup.emptyBackground': 'var(--code-base)',
    'editorGroupHeader.tabsBackground': 'var(--code-crust)',
    'editorGutter.addedBackground': 'var(--code-green)',
    'editorGutter.background': 'var(--code-base)',
    'editorGutter.commentRangeForeground': 'var(--code-surface-0)',
    'editorGutter.commentGlyphForeground': 'var(--code-mauve)',
    'editorGutter.deletedBackground': 'var(--code-red)',
    'editorGutter.foldingControlForeground': 'var(--code-overlay-2)',
    'editorGutter.modifiedBackground': 'var(--code-yellow)',
    'editorHoverWidget.background': 'var(--code-mantle)',
    'editorHoverWidget.border': 'var(--code-surface-2)',
    'editorHoverWidget.foreground': 'var(--code-text)',
    'editorIndentGuide.activeBackground': 'var(--code-surface-2)',
    'editorIndentGuide.background': 'var(--code-surface-1)',
    'editorInlayHint.foreground': 'var(--code-surface-2)',
    'editorInlayHint.background': 'var(--code-mantle)bf',
    'editorInlayHint.typeForeground': 'var(--code-subtext-1)',
    'editorInlayHint.typeBackground': 'var(--code-mantle)bf',
    'editorInlayHint.parameterForeground': 'var(--code-subtext-0)',
    'editorInlayHint.parameterBackground': 'var(--code-mantle)bf',
    'editorLineNumber.activeForeground': 'var(--code-mauve)',
    'editorLineNumber.foreground': 'var(--code-overlay-1)',
    'editorLink.activeForeground': 'var(--code-mauve)',
    'editorMarkerNavigation.background': 'var(--code-mantle)',
    'editorMarkerNavigationError.background': 'var(--code-red)',
    'editorMarkerNavigationInfo.background': 'var(--code-blue)',
    'editorMarkerNavigationWarning.background': 'var(--code-peach)',
    'editorOverviewRuler.background': 'var(--code-mantle)',
    'editorOverviewRuler.border': 'var(--code-text)12',
    'editorOverviewRuler.modifiedForeground': 'var(--code-yellow)',
    'editorRuler.foreground': 'var(--code-surface-2)',
    'editor.stackFrameHighlightBackground': 'var(--code-yellow)26',
    'editor.focusedStackFrameHighlightBackground': 'var(--code-green)26',
    'editorStickyScrollHover.background': 'var(--code-surface-0)',
    'editorSuggestWidget.background': 'var(--code-mantle)',
    'editorSuggestWidget.border': 'var(--code-surface-2)',
    'editorSuggestWidget.foreground': 'var(--code-text)',
    'editorSuggestWidget.highlightForeground': 'var(--code-mauve)',
    'editorSuggestWidget.selectedBackground': 'var(--code-surface-0)',
    'editorWhitespace.foreground': 'var(--code-overlay-2)66',
    'editorWidget.background': 'var(--code-mantle)',
    'editorWidget.foreground': 'var(--code-text)',
    'editorWidget.resizeBorder': 'var(--code-surface-2)',
    'editorLightBulb.foreground': 'var(--code-yellow)',
    'editorError.foreground': 'var(--code-red)',
    'editorError.border': 'var(--code-black)',
    'editorError.background': 'var(--code-black)',
    'editorWarning.foreground': 'var(--code-peach)',
    'editorWarning.border': 'var(--code-black)',
    'editorWarning.background': 'var(--code-black)',
    'editorInfo.foreground': 'var(--code-blue)',
    'editorInfo.border': 'var(--code-black)',
    'editorInfo.background': 'var(--code-black)',
    'problemsErrorIcon.foreground': 'var(--code-red)',
    'problemsInfoIcon.foreground': 'var(--code-blue)',
    'problemsWarningIcon.foreground': 'var(--code-peach)',
    'extensionButton.prominentForeground': 'var(--code-crust)',
    'extensionButton.prominentBackground': 'var(--code-mauve)',
    'extensionButton.separator': 'var(--code-base)',
    'extensionButton.prominentHoverBackground': '#dac1f9',
    'extensionBadge.remoteBackground': 'var(--code-blue)',
    'extensionBadge.remoteForeground': 'var(--code-crust)',
    'extensionIcon.starForeground': 'var(--code-yellow)',
    'extensionIcon.verifiedForeground': 'var(--code-green)',
    'extensionIcon.preReleaseForeground': 'var(--code-surface-2)',
    'extensionIcon.sponsorForeground': 'var(--code-pink)',
    'gitDecoration.addedResourceForeground': 'var(--code-green)',
    'gitDecoration.conflictingResourceForeground': 'var(--code-mauve)',
    'gitDecoration.deletedResourceForeground': 'var(--code-red)',
    'gitDecoration.ignoredResourceForeground': 'var(--code-overlay-0)',
    'gitDecoration.modifiedResourceForeground': 'var(--code-yellow)',
    'gitDecoration.stageDeletedResourceForeground': 'var(--code-red)',
    'gitDecoration.stageModifiedResourceForeground': 'var(--code-yellow)',
    'gitDecoration.submoduleResourceForeground': 'var(--code-blue)',
    'gitDecoration.untrackedResourceForeground': 'var(--code-green)',
    'input.background': 'var(--code-surface-0)',
    'input.border': 'var(--code-black)',
    'input.foreground': 'var(--code-text)',
    'input.placeholderForeground': 'var(--code-text)73',
    'inputOption.activeBackground': 'var(--code-surface-2)',
    'inputOption.activeBorder': 'var(--code-mauve)',
    'inputOption.activeForeground': 'var(--code-text)',
    'inputValidation.errorBackground': 'var(--code-red)',
    'inputValidation.errorBorder': 'var(--code-crust)33',
    'inputValidation.errorForeground': 'var(--code-crust)',
    'inputValidation.infoBackground': 'var(--code-blue)',
    'inputValidation.infoBorder': 'var(--code-crust)33',
    'inputValidation.infoForeground': 'var(--code-crust)',
    'inputValidation.warningBackground': 'var(--code-peach)',
    'inputValidation.warningBorder': 'var(--code-crust)33',
    'inputValidation.warningForeground': 'var(--code-crust)',
    'list.activeSelectionBackground': 'var(--code-surface-0)',
    'list.activeSelectionForeground': 'var(--code-text)',
    'list.dropBackground': 'var(--code-mauve)33',
    'list.focusBackground': 'var(--code-surface-0)',
    'list.focusForeground': 'var(--code-text)',
    'list.focusOutline': 'var(--code-black)',
    'list.highlightForeground': 'var(--code-mauve)',
    'list.hoverBackground': 'var(--code-surface-0)80',
    'list.hoverForeground': 'var(--code-text)',
    'list.inactiveSelectionBackground': 'var(--code-surface-0)',
    'list.inactiveSelectionForeground': 'var(--code-text)',
    'list.warningForeground': 'var(--code-peach)',
    'listFilterWidget.background': 'var(--code-surface-1)',
    'listFilterWidget.noMatchesOutline': 'var(--code-red)',
    'listFilterWidget.outline': 'var(--code-black)',
    'tree.indentGuidesStroke': 'var(--code-overlay-2)',
    'tree.inactiveIndentGuidesStroke': 'var(--code-surface-1)',
    'menu.background': 'var(--code-base)',
    'menu.border': 'var(--code-base)80',
    'menu.foreground': 'var(--code-text)',
    'menu.selectionBackground': 'var(--code-surface-2)',
    'menu.selectionBorder': 'var(--code-black)',
    'menu.selectionForeground': 'var(--code-text)',
    'menu.separatorBackground': 'var(--code-surface-2)',
    'menubar.selectionBackground': 'var(--code-surface-1)',
    'menubar.selectionForeground': 'var(--code-text)',
    'merge.commonContentBackground': 'var(--code-surface-1)',
    'merge.commonHeaderBackground': 'var(--code-surface-2)',
    'merge.currentContentBackground': 'var(--code-green)33',
    'merge.currentHeaderBackground': 'var(--code-green)66',
    'merge.incomingContentBackground': 'var(--code-blue)33',
    'merge.incomingHeaderBackground': 'var(--code-blue)66',
    'minimap.background': 'var(--code-mantle)80',
    'minimap.findMatchHighlight': 'var(--code-sky)4d',
    'minimap.selectionHighlight': 'var(--code-surface-2)bf',
    'minimap.selectionOccurrenceHighlight': 'var(--code-surface-2)bf',
    'minimap.warningHighlight': 'var(--code-peach)bf',
    'minimap.errorHighlight': 'var(--code-red)bf',
    'minimapSlider.background': 'var(--code-mauve)33',
    'minimapSlider.hoverBackground': 'var(--code-mauve)66',
    'minimapSlider.activeBackground': 'var(--code-mauve)99',
    'minimapGutter.addedBackground': 'var(--code-green)bf',
    'minimapGutter.deletedBackground': 'var(--code-red)bf',
    'minimapGutter.modifiedBackground': 'var(--code-yellow)bf',
    'notificationCenter.border': 'var(--code-mauve)',
    'notificationCenterHeader.foreground': 'var(--code-text)',
    'notificationCenterHeader.background': 'var(--code-mantle)',
    'notificationToast.border': 'var(--code-mauve)',
    'notifications.foreground': 'var(--code-text)',
    'notifications.background': 'var(--code-mantle)',
    'notifications.border': 'var(--code-mauve)',
    'notificationLink.foreground': 'var(--code-blue)',
    'notificationsErrorIcon.foreground': 'var(--code-red)',
    'notificationsWarningIcon.foreground': 'var(--code-peach)',
    'notificationsInfoIcon.foreground': 'var(--code-blue)',
    'panel.background': 'var(--code-base)',
    'panel.border': 'var(--code-surface-2)',
    'panelSection.border': 'var(--code-surface-2)',
    'panelSection.dropBackground': 'var(--code-mauve)33',
    'panelTitle.activeBorder': 'var(--code-mauve)',
    'panelTitle.activeForeground': 'var(--code-text)',
    'panelTitle.inactiveForeground': 'var(--code-subtext-0)',
    'peekView.border': 'var(--code-mauve)',
    'peekViewEditor.background': 'var(--code-mantle)',
    'peekViewEditorGutter.background': 'var(--code-mantle)',
    'peekViewEditor.matchHighlightBackground': 'var(--code-sky)4d',
    'peekViewEditor.matchHighlightBorder': 'var(--code-black)',
    'peekViewResult.background': 'var(--code-mantle)',
    'peekViewResult.fileForeground': 'var(--code-text)',
    'peekViewResult.lineForeground': 'var(--code-text)',
    'peekViewResult.matchHighlightBackground': 'var(--code-sky)4d',
    'peekViewResult.selectionBackground': 'var(--code-surface-0)',
    'peekViewResult.selectionForeground': 'var(--code-text)',
    'peekViewTitle.background': 'var(--code-base)',
    'peekViewTitleDescription.foreground': 'var(--code-subtext-1)b3',
    'peekViewTitleLabel.foreground': 'var(--code-text)',
    'pickerGroup.border': 'var(--code-mauve)',
    'pickerGroup.foreground': 'var(--code-mauve)',
    'progressBar.background': 'var(--code-mauve)',
    'scrollbar.shadow': 'var(--code-crust)',
    'scrollbarSlider.activeBackground': 'var(--code-surface-0)66',
    'scrollbarSlider.background': 'var(--code-surface-2)80',
    'scrollbarSlider.hoverBackground': 'var(--code-overlay-0)',
    'settings.focusedRowBackground': 'var(--code-surface-2)33',
    'settings.headerForeground': 'var(--code-text)',
    'settings.modifiedItemIndicator': 'var(--code-mauve)',
    'settings.dropdownBackground': 'var(--code-surface-1)',
    'settings.dropdownListBorder': 'var(--code-black)',
    'settings.textInputBackground': 'var(--code-surface-1)',
    'settings.textInputBorder': 'var(--code-black)',
    'settings.numberInputBackground': 'var(--code-surface-1)',
    'settings.numberInputBorder': 'var(--code-black)',
    'sideBar.background': 'var(--code-mantle)',
    'sideBar.dropBackground': 'var(--code-mauve)33',
    'sideBar.foreground': 'var(--code-text)',
    'sideBar.border': 'var(--code-black)',
    'sideBarSectionHeader.background': 'var(--code-mantle)',
    'sideBarSectionHeader.foreground': 'var(--code-text)',
    'sideBarTitle.foreground': 'var(--code-mauve)',
    'banner.background': 'var(--code-surface-1)',
    'banner.foreground': 'var(--code-text)',
    'banner.iconForeground': 'var(--code-text)',
    'statusBar.background': 'var(--code-crust)',
    'statusBar.foreground': 'var(--code-text)',
    'statusBar.border': 'var(--code-black)',
    'statusBar.noFolderBackground': 'var(--code-crust)',
    'statusBar.noFolderForeground': 'var(--code-text)',
    'statusBar.noFolderBorder': 'var(--code-black)',
    'statusBar.debuggingBackground': 'var(--code-peach)',
    'statusBar.debuggingForeground': 'var(--code-crust)',
    'statusBar.debuggingBorder': 'var(--code-black)',
    'statusBarItem.remoteBackground': 'var(--code-blue)',
    'statusBarItem.remoteForeground': 'var(--code-crust)',
    'statusBarItem.activeBackground': 'var(--code-surface-2)66',
    'statusBarItem.hoverBackground': 'var(--code-surface-2)33',
    'statusBarItem.prominentForeground': 'var(--code-mauve)',
    'statusBarItem.prominentBackground': 'var(--code-black)',
    'statusBarItem.prominentHoverBackground': 'var(--code-surface-2)33',
    'statusBarItem.errorForeground': 'var(--code-red)',
    'statusBarItem.errorBackground': 'var(--code-black)',
    'statusBarItem.warningForeground': 'var(--code-peach)',
    'statusBarItem.warningBackground': 'var(--code-black)',
    'commandCenter.foreground': 'var(--code-subtext-1)',
    'commandCenter.inactiveForeground': 'var(--code-subtext-1)',
    'commandCenter.activeForeground': 'var(--code-mauve)',
    'commandCenter.background': 'var(--code-mantle)',
    'commandCenter.activeBackground': 'var(--code-surface-2)33',
    'commandCenter.border': 'var(--code-black)',
    'commandCenter.inactiveBorder': 'var(--code-black)',
    'commandCenter.activeBorder': 'var(--code-mauve)',
    'tab.activeBackground': 'var(--code-base)',
    'tab.activeBorder': 'var(--code-black)',
    'tab.activeBorderTop': 'var(--code-mauve)',
    'tab.activeForeground': 'var(--code-mauve)',
    'tab.activeModifiedBorder': 'var(--code-yellow)',
    'tab.border': 'var(--code-mantle)',
    'tab.hoverBackground': '#2e324a',
    'tab.hoverBorder': 'var(--code-black)',
    'tab.hoverForeground': 'var(--code-mauve)',
    'tab.inactiveBackground': 'var(--code-mantle)',
    'tab.inactiveForeground': 'var(--code-overlay-0)',
    'tab.inactiveModifiedBorder': 'var(--code-yellow)4d',
    'tab.lastPinnedBorder': 'var(--code-mauve)',
    'tab.unfocusedActiveBackground': 'var(--code-mantle)',
    'tab.unfocusedActiveBorder': 'var(--code-black)',
    'tab.unfocusedActiveBorderTop': 'var(--code-mauve)4d',
    'tab.unfocusedInactiveBackground': '#141620',
    'terminal.foreground': 'var(--code-text)',
    'terminal.ansiBlack': 'var(--code-subtext-0)',
    'terminal.ansiRed': 'var(--code-red)',
    'terminal.ansiGreen': 'var(--code-green)',
    'terminal.ansiYellow': 'var(--code-yellow)',
    'terminal.ansiBlue': 'var(--code-blue)',
    'terminal.ansiMagenta': 'var(--code-pink)',
    'terminal.ansiCyan': 'var(--code-sky)',
    'terminal.ansiWhite': 'var(--code-subtext-1)',
    'terminal.ansiBrightBlack': 'var(--code-surface-2)',
    'terminal.ansiBrightRed': 'var(--code-red)',
    'terminal.ansiBrightGreen': 'var(--code-green)',
    'terminal.ansiBrightYellow': 'var(--code-yellow)',
    'terminal.ansiBrightBlue': 'var(--code-blue)',
    'terminal.ansiBrightMagenta': 'var(--code-pink)',
    'terminal.ansiBrightCyan': 'var(--code-sky)',
    'terminal.ansiBrightWhite': 'var(--code-surface-1)',
    'terminal.selectionBackground': 'var(--code-surface-2)',
    'terminal.inactiveSelectionBackground': 'var(--code-surface-2)80',
    'terminalCursor.background': 'var(--code-base)',
    'terminalCursor.foreground': 'var(--code-rosewater)',
    'terminal.border': 'var(--code-surface-2)',
    'terminal.dropBackground': 'var(--code-mauve)33',
    'terminal.tab.activeBorder': 'var(--code-mauve)',
    'terminalCommandDecoration.defaultBackground': 'var(--code-surface-2)',
    'terminalCommandDecoration.successBackground': 'var(--code-green)',
    'terminalCommandDecoration.errorBackground': 'var(--code-red)',
    'titleBar.activeBackground': 'var(--code-crust)',
    'titleBar.activeForeground': 'var(--code-text)',
    'titleBar.inactiveBackground': 'var(--code-crust)',
    'titleBar.inactiveForeground': 'var(--code-text)80',
    'titleBar.border': 'var(--code-black)',
    'welcomePage.tileBackground': 'var(--code-mantle)',
    'welcomePage.progress.background': 'var(--code-crust)',
    'welcomePage.progress.foreground': 'var(--code-mauve)',
    'walkThrough.embeddedEditorBackground': 'var(--code-base)4d',
    'symbolIcon.textForeground': 'var(--code-text)',
    'symbolIcon.arrayForeground': 'var(--code-peach)',
    'symbolIcon.booleanForeground': 'var(--code-mauve)',
    'symbolIcon.classForeground': 'var(--code-yellow)',
    'symbolIcon.colorForeground': 'var(--code-pink)',
    'symbolIcon.constantForeground': 'var(--code-peach)',
    'symbolIcon.constructorForeground': 'var(--code-lavender)',
    'symbolIcon.enumeratorForeground': 'var(--code-yellow)',
    'symbolIcon.enumeratorMemberForeground': 'var(--code-yellow)',
    'symbolIcon.eventForeground': 'var(--code-pink)',
    'symbolIcon.fieldForeground': 'var(--code-text)',
    'symbolIcon.fileForeground': 'var(--code-mauve)',
    'symbolIcon.folderForeground': 'var(--code-mauve)',
    'symbolIcon.functionForeground': 'var(--code-blue)',
    'symbolIcon.interfaceForeground': 'var(--code-yellow)',
    'symbolIcon.keyForeground': 'var(--code-teal)',
    'symbolIcon.keywordForeground': 'var(--code-mauve)',
    'symbolIcon.methodForeground': 'var(--code-blue)',
    'symbolIcon.moduleForeground': 'var(--code-text)',
    'symbolIcon.namespaceForeground': 'var(--code-yellow)',
    'symbolIcon.nullForeground': 'var(--code-maroon)',
    'symbolIcon.numberForeground': 'var(--code-peach)',
    'symbolIcon.objectForeground': 'var(--code-yellow)',
    'symbolIcon.operatorForeground': 'var(--code-teal)',
    'symbolIcon.packageForeground': 'var(--code-flamingo)',
    'symbolIcon.propertyForeground': 'var(--code-maroon)',
    'symbolIcon.referenceForeground': 'var(--code-yellow)',
    'symbolIcon.snippetForeground': 'var(--code-flamingo)',
    'symbolIcon.stringForeground': 'var(--code-green)',
    'symbolIcon.structForeground': 'var(--code-teal)',
    'symbolIcon.typeParameterForeground': 'var(--code-maroon)',
    'symbolIcon.unitForeground': 'var(--code-text)',
    'symbolIcon.variableForeground': 'var(--code-text)',
    'charts.foreground': 'var(--code-text)',
    'charts.lines': 'var(--code-subtext-1)',
    'charts.red': 'var(--code-red)',
    'charts.blue': 'var(--code-blue)',
    'charts.yellow': 'var(--code-yellow)',
    'charts.orange': 'var(--code-peach)',
    'charts.green': 'var(--code-green)',
    'charts.purple': 'var(--code-mauve)',
    'errorLens.errorBackground': 'var(--code-red)26',
    'errorLens.errorBackgroundLight': 'var(--code-red)26',
    'errorLens.errorForeground': 'var(--code-red)',
    'errorLens.errorForegroundLight': 'var(--code-red)',
    'errorLens.errorMessageBackground': 'var(--code-red)26',
    'errorLens.hintBackground': 'var(--code-green)26',
    'errorLens.hintBackgroundLight': 'var(--code-green)26',
    'errorLens.hintForeground': 'var(--code-green)',
    'errorLens.hintForegroundLight': 'var(--code-green)',
    'errorLens.hintMessageBackground': 'var(--code-green)26',
    'errorLens.infoBackground': 'var(--code-blue)26',
    'errorLens.infoBackgroundLight': 'var(--code-blue)26',
    'errorLens.infoForeground': 'var(--code-blue)',
    'errorLens.infoForegroundLight': 'var(--code-blue)',
    'errorLens.infoMessageBackground': 'var(--code-blue)26',
    'errorLens.statusBarErrorForeground': 'var(--code-red)',
    'errorLens.statusBarHintForeground': 'var(--code-green)',
    'errorLens.statusBarIconErrorForeground': 'var(--code-red)',
    'errorLens.statusBarIconWarningForeground': 'var(--code-peach)',
    'errorLens.statusBarInfoForeground': 'var(--code-blue)',
    'errorLens.statusBarWarningForeground': 'var(--code-peach)',
    'errorLens.warningBackground': 'var(--code-peach)26',
    'errorLens.warningBackgroundLight': 'var(--code-peach)26',
    'errorLens.warningForeground': 'var(--code-peach)',
    'errorLens.warningForegroundLight': 'var(--code-peach)',
    'errorLens.warningMessageBackground': 'var(--code-peach)26',
    'issues.closed': 'var(--code-mauve)',
    'issues.newIssueDecoration': 'var(--code-rosewater)',
    'issues.open': 'var(--code-green)',
    'pullRequests.closed': 'var(--code-red)',
    'pullRequests.draft': 'var(--code-overlay-2)',
    'pullRequests.merged': 'var(--code-mauve)',
    'pullRequests.notification': 'var(--code-text)',
    'pullRequests.open': 'var(--code-green)',
    'gitlens.gutterBackgroundColor': 'var(--code-surface-0)4d',
    'gitlens.gutterForegroundColor': 'var(--code-text)',
    'gitlens.gutterUncommittedForegroundColor': 'var(--code-mauve)',
    'gitlens.trailingLineBackgroundColor': 'var(--code-black)',
    'gitlens.trailingLineForegroundColor': 'var(--code-text)4d',
    'gitlens.lineHighlightBackgroundColor': 'var(--code-mauve)26',
    'gitlens.lineHighlightOverviewRulerColor': 'var(--code-mauve)cc',
    'gitlens.openAutolinkedIssueIconColor': 'var(--code-green)',
    'gitlens.closedAutolinkedIssueIconColor': 'var(--code-mauve)',
    'gitlens.closedPullRequestIconColor': 'var(--code-red)',
    'gitlens.openPullRequestIconColor': 'var(--code-green)',
    'gitlens.mergedPullRequestIconColor': 'var(--code-mauve)',
    'gitlens.unpublishedChangesIconColor': 'var(--code-green)',
    'gitlens.unpublishedCommitIconColor': 'var(--code-green)',
    'gitlens.unpulledChangesIconColor': 'var(--code-peach)',
    'gitlens.decorations.branchAheadForegroundColor': 'var(--code-green)',
    'gitlens.decorations.branchBehindForegroundColor': 'var(--code-peach)',
    'gitlens.decorations.branchDivergedForegroundColor': 'var(--code-yellow)',
    'gitlens.decorations.branchUnpublishedForegroundColor': 'var(--code-green)',
    'gitlens.decorations.branchMissingUpstreamForegroundColor':
      'var(--code-peach)',
    'gitlens.decorations.statusMergingOrRebasingConflictForegroundColor':
      'var(--code-maroon)',
    'gitlens.decorations.statusMergingOrRebasingForegroundColor':
      'var(--code-yellow)',
    'gitlens.decorations.workspaceRepoMissingForegroundColor':
      'var(--code-subtext-0)',
    'gitlens.decorations.workspaceCurrentForegroundColor': 'var(--code-mauve)',
    'gitlens.decorations.workspaceRepoOpenForegroundColor': 'var(--code-mauve)',
    'gitlens.decorations.worktreeHasUncommittedChangesForegroundColor':
      'var(--code-peach)',
    'gitlens.decorations.worktreeMissingForegroundColor': 'var(--code-maroon)',
    'gitlens.graphLane1Color': 'var(--code-mauve)',
    'gitlens.graphLane2Color': 'var(--code-yellow)',
    'gitlens.graphLane3Color': 'var(--code-blue)',
    'gitlens.graphLane4Color': 'var(--code-flamingo)',
    'gitlens.graphLane5Color': 'var(--code-green)',
    'gitlens.graphLane6Color': 'var(--code-lavender)',
    'gitlens.graphLane7Color': 'var(--code-rosewater)',
    'gitlens.graphLane8Color': 'var(--code-red)',
    'gitlens.graphLane9Color': 'var(--code-teal)',
    'gitlens.graphLane10Color': 'var(--code-pink)',
    'gitlens.graphChangesColumnAddedColor': 'var(--code-green)',
    'gitlens.graphChangesColumnDeletedColor': 'var(--code-red)',
    'gitlens.graphMinimapMarkerHeadColor': 'var(--code-green)',
    'gitlens.graphScrollMarkerHeadColor': 'var(--code-green)',
    'gitlens.graphMinimapMarkerUpstreamColor': '#96d382',
    'gitlens.graphScrollMarkerUpstreamColor': '#96d382',
    'gitlens.graphMinimapMarkerHighlightsColor': 'var(--code-yellow)',
    'gitlens.graphScrollMarkerHighlightsColor': 'var(--code-yellow)',
    'gitlens.graphMinimapMarkerLocalBranchesColor': 'var(--code-blue)',
    'gitlens.graphScrollMarkerLocalBranchesColor': 'var(--code-blue)',
    'gitlens.graphMinimapMarkerRemoteBranchesColor': '#739df2',
    'gitlens.graphScrollMarkerRemoteBranchesColor': '#739df2',
    'gitlens.graphMinimapMarkerStashesColor': 'var(--code-mauve)',
    'gitlens.graphScrollMarkerStashesColor': 'var(--code-mauve)',
    'gitlens.graphMinimapMarkerTagsColor': 'var(--code-flamingo)',
    'gitlens.graphScrollMarkerTagsColor': 'var(--code-flamingo)',
    'editorBracketHighlight.foreground1': 'var(--code-red)',
    'editorBracketHighlight.foreground2': 'var(--code-peach)',
    'editorBracketHighlight.foreground3': 'var(--code-yellow)',
    'editorBracketHighlight.foreground4': 'var(--code-green)',
    'editorBracketHighlight.foreground5': 'var(--code-sapphire)',
    'editorBracketHighlight.foreground6': 'var(--code-mauve)',
    'editorBracketHighlight.unexpectedBracket.foreground': 'var(--code-maroon)',
    'button.secondaryBorder': 'var(--code-mauve)',
    'table.headerBackground': 'var(--code-surface-0)',
    'table.headerForeground': 'var(--code-text)',
    'list.focusAndSelectionBackground': 'var(--code-surface-1)',
  },
  semanticHighlighting: true,
  semanticTokenColors: {
    enumMember: {
      foreground: 'var(--code-teal)',
    },
    selfKeyword: {
      foreground: 'var(--code-red)',
    },
    boolean: {
      foreground: 'var(--code-peach)',
    },
    number: {
      foreground: 'var(--code-peach)',
    },
    'variable.defaultLibrary': {
      foreground: 'var(--code-maroon)',
    },
    'class:python': {
      foreground: 'var(--code-yellow)',
    },
    'class.builtin:python': {
      foreground: 'var(--code-mauve)',
    },
    'variable.typeHint:python': {
      foreground: 'var(--code-yellow)',
    },
    'function.decorator:python': {
      foreground: 'var(--code-peach)',
    },
    'variable.readonly:javascript': {
      foreground: 'var(--code-text)',
    },
    'variable.readonly:typescript': {
      foreground: 'var(--code-text)',
    },
    'property.readonly:javascript': {
      foreground: 'var(--code-text)',
    },
    'property.readonly:typescript': {
      foreground: 'var(--code-text)',
    },
    'variable.readonly:javascriptreact': {
      foreground: 'var(--code-text)',
    },
    'variable.readonly:typescriptreact': {
      foreground: 'var(--code-text)',
    },
    'property.readonly:javascriptreact': {
      foreground: 'var(--code-text)',
    },
    'property.readonly:typescriptreact': {
      foreground: 'var(--code-text)',
    },
    'variable.readonly:scala': {
      foreground: 'var(--code-text)',
    },
    'type.defaultLibrary:go': {
      foreground: 'var(--code-mauve)',
    },
    'variable.readonly.defaultLibrary:go': {
      foreground: 'var(--code-mauve)',
    },
    tomlArrayKey: {
      foreground: 'var(--code-blue)',
      fontStyle: '',
    },
    tomlTableKey: {
      foreground: 'var(--code-blue)',
      fontStyle: '',
    },
    'builtinAttribute.attribute.library:rust': {
      foreground: 'var(--code-blue)',
    },
    'generic.attribute:rust': {
      foreground: 'var(--code-text)',
    },
    'constant.builtin.readonly:nix': {
      foreground: 'var(--code-mauve)',
    },
    heading: {
      foreground: 'var(--code-red)',
    },
    'text.emph': {
      foreground: 'var(--code-red)',
      fontStyle: 'italic',
    },
    'text.strong': {
      foreground: 'var(--code-red)',
      fontStyle: 'bold',
    },
    'text.math': {
      foreground: 'var(--code-flamingo)',
    },
    pol: {
      foreground: 'var(--code-flamingo)',
    },
  },
  tokenColors: [
    {
      name: 'Basic text & variable names (incl. leading punctuation)',
      scope: [
        'text',
        'source',
        'variable.other.readwrite',
        'punctuation.definition.variable',
      ],
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      name: 'Parentheses, Brackets, Braces',
      scope: 'punctuation',
      settings: {
        foreground: 'var(--code-overlay-2)',
        fontStyle: '',
      },
    },
    {
      name: 'Comments',
      scope: ['comment', 'punctuation.definition.comment'],
      settings: {
        foreground: 'var(--code-overlay-0)',
        fontStyle: 'italic',
      },
    },
    {
      scope: ['string', 'punctuation.definition.string'],
      settings: {
        foreground: 'var(--code-green)',
      },
    },
    {
      scope: 'constant.character.escape',
      settings: {
        foreground: 'var(--code-pink)',
      },
    },
    {
      name: 'Booleans, constants, numbers',
      scope: [
        'constant.numeric',
        'variable.other.constant',
        'entity.name.constant',
        'constant.language.boolean',
        'constant.language.false',
        'constant.language.true',
        'keyword.other.unit.user-defined',
        'keyword.other.unit.suffix.floating-point',
      ],
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      scope: [
        'keyword',
        'keyword.operator.word',
        'keyword.operator.new',
        'variable.language.super',
        'support.type.primitive',
        'storage.type',
        'storage.modifier',
        'punctuation.definition.keyword',
      ],
      settings: {
        foreground: 'var(--code-keyword)',
        fontStyle: 'normal',
      },
    },
    {
      scope: 'entity.name.tag.documentation',
      settings: {
        foreground: 'var(--code-mauve)',
      },
    },
    {
      name: 'Punctuation',
      scope: [
        'keyword.operator',
        'punctuation.accessor',
        'punctuation.definition.generic',
        'meta.function.closure punctuation.section.parameters',
        'punctuation.definition.tag',
        'punctuation.separator.key-value',
      ],
      settings: {
        foreground: 'var(--code-punctuation)',
        fontStyle: 'normal',
      },
    },
    {
      scope: [
        'entity.name.function',
        'meta.function-call.method',
        'support.function',
        'support.function.misc',
        'variable.function',
      ],
      settings: {
        foreground: 'var(--code-function)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Classes',
      scope: [
        'entity.name.class',
        'entity.other.inherited-class',
        'support.class',
        'meta.function-call.constructor',
        'entity.name.struct',
      ],
      settings: {
        foreground: 'var(--code-yellow)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Enum',
      scope: 'entity.name.enum',
      settings: {
        foreground: 'var(--code-yellow)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Enum member',
      scope: [
        'meta.enum variable.other.readwrite',
        'variable.other.enummember',
      ],
      settings: {
        foreground: 'var(--code-teal)',
      },
    },
    {
      name: 'Object properties',
      scope: 'meta.property.object',
      settings: {
        foreground: 'var(--code-teal)',
      },
    },
    {
      name: 'Types',
      scope: [
        'meta.type',
        'meta.type-alias',
        'support.type',
        'entity.name.type',
      ],
      settings: {
        foreground: 'var(--code-yellow)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Decorators',
      scope: [
        'meta.annotation variable.function',
        'meta.annotation variable.annotation.function',
        'meta.annotation punctuation.definition.annotation',
        'meta.decorator',
        'punctuation.decorator',
      ],
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      scope: ['variable.parameter', 'meta.function.parameters'],
      settings: {
        foreground: 'var(--code-variable)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Built-ins',
      scope: ['constant.language', 'support.function.builtin'],
      settings: {
        foreground: 'var(--code-red)',
      },
    },
    {
      scope: 'entity.other.attribute-name.documentation',
      settings: {
        foreground: 'var(--code-red)',
      },
    },
    {
      name: 'Preprocessor directives',
      scope: ['keyword.control.directive', 'punctuation.definition.directive'],
      settings: {
        foreground: 'var(--code-yellow)',
      },
    },
    {
      name: 'Type parameters',
      scope: 'punctuation.definition.typeparameters',
      settings: {
        foreground: 'var(--code-sky)',
      },
    },
    {
      name: 'Namespaces',
      scope: 'entity.name.namespace',
      settings: {
        foreground: 'var(--code-yellow)',
      },
    },
    {
      name: 'Property names (left hand assignments in json/yaml/css)',
      scope: 'support.type.property-name.css',
      settings: {
        foreground: 'var(--code-blue)',
        fontStyle: '',
      },
    },
    {
      name: 'This/Self keyword',
      scope: [
        'variable.language.this',
        'variable.language.this punctuation.definition.variable',
      ],
      settings: {
        foreground: 'var(--code-red)',
      },
    },
    {
      name: 'Object properties',
      scope: 'variable.object.property',
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      name: 'String template interpolation',
      scope: ['string.template variable', 'string variable'],
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      name: '`new` as bold',
      scope: 'keyword.operator.new',
      settings: {
        fontStyle: 'bold',
      },
    },
    {
      name: 'C++ extern keyword',
      scope: 'storage.modifier.specifier.extern.cpp',
      settings: {
        foreground: 'var(--code-mauve)',
      },
    },
    {
      name: 'C++ scope resolution',
      scope: [
        'entity.name.scope-resolution.template.call.cpp',
        'entity.name.scope-resolution.parameter.cpp',
        'entity.name.scope-resolution.cpp',
        'entity.name.scope-resolution.function.definition.cpp',
      ],
      settings: {
        foreground: 'var(--code-yellow)',
      },
    },
    {
      name: 'C++ doc keywords',
      scope: 'storage.type.class.doxygen',
      settings: {
        fontStyle: '',
      },
    },
    {
      name: 'C++ operators',
      scope: ['storage.modifier.reference.cpp'],
      settings: {
        foreground: 'var(--code-teal)',
      },
    },
    {
      name: 'C# Interpolated Strings',
      scope: 'meta.interpolation.cs',
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      name: 'C# xml-style docs',
      scope: 'comment.block.documentation.cs',
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      name: 'Classes, reflecting the className color in JSX',
      scope: [
        'source.css entity.other.attribute-name.class.css',
        'entity.other.attribute-name.parent-selector.css punctuation.definition.entity.css',
      ],
      settings: {
        foreground: 'var(--code-yellow)',
      },
    },
    {
      name: 'Operators',
      scope: 'punctuation.separator.operator.css',
      settings: {
        foreground: 'var(--code-teal)',
      },
    },
    {
      name: 'Pseudo classes',
      scope: 'source.css entity.other.attribute-name.pseudo-class',
      settings: {
        foreground: 'var(--code-teal)',
      },
    },
    {
      scope: 'source.css constant.other.unicode-range',
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      scope: 'source.css variable.parameter.url',
      settings: {
        foreground: 'var(--code-green)',
        fontStyle: '',
      },
    },
    {
      name: 'CSS vendored property names',
      scope: ['support.type.vendored.property-name'],
      settings: {
        foreground: 'var(--code-sky)',
      },
    },
    {
      name: 'Less/SCSS right-hand variables (@/$-prefixed)',
      scope: [
        'source.css meta.property-value variable',
        'source.css meta.property-value variable.other.less',
        'source.css meta.property-value variable.other.less punctuation.definition.variable.less',
        'meta.definition.variable.scss',
      ],
      settings: {
        foreground: 'var(--code-maroon)',
      },
    },
    {
      name: 'CSS variables (--prefixed)',
      scope: [
        'source.css meta.property-list variable',
        'meta.property-list variable.other.less',
        'meta.property-list variable.other.less punctuation.definition.variable.less',
      ],
      settings: {
        foreground: 'var(--code-blue)',
      },
    },
    {
      name: 'CSS Percentage values, styled the same as numbers',
      scope: 'keyword.other.unit.percentage.css',
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      name: 'CSS Attribute selectors, styled the same as strings',
      scope: 'source.css meta.attribute-selector',
      settings: {
        foreground: 'var(--code-green)',
      },
    },
    {
      name: 'JSON/YAML keys, other left-hand assignments',
      scope: [
        'keyword.other.definition.ini',
        'punctuation.support.type.property-name.json',
        'support.type.property-name.json',
        'punctuation.support.type.property-name.toml',
        'support.type.property-name.toml',
        'entity.name.tag.yaml',
        'punctuation.support.type.property-name.yaml',
        'support.type.property-name.yaml',
      ],
      settings: {
        foreground: 'var(--code-blue)',
        fontStyle: '',
      },
    },
    {
      name: 'JSON/YAML constants',
      scope: ['constant.language.json', 'constant.language.yaml'],
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      name: 'YAML anchors',
      scope: ['entity.name.type.anchor.yaml', 'variable.other.alias.yaml'],
      settings: {
        foreground: 'var(--code-yellow)',
        fontStyle: '',
      },
    },
    {
      name: 'TOML tables / ini groups',
      scope: [
        'support.type.property-name.table',
        'entity.name.section.group-title.ini',
      ],
      settings: {
        foreground: 'var(--code-yellow)',
      },
    },
    {
      name: 'TOML dates',
      scope: 'constant.other.time.datetime.offset.toml',
      settings: {
        foreground: 'var(--code-pink)',
      },
    },
    {
      name: 'YAML anchor puctuation',
      scope: [
        'punctuation.definition.anchor.yaml',
        'punctuation.definition.alias.yaml',
      ],
      settings: {
        foreground: 'var(--code-pink)',
      },
    },
    {
      name: 'YAML triple dashes',
      scope: 'entity.other.document.begin.yaml',
      settings: {
        foreground: 'var(--code-pink)',
      },
    },
    {
      name: 'Markup Diff',
      scope: 'markup.changed.diff',
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      name: 'Diff',
      scope: [
        'meta.diff.header.from-file',
        'meta.diff.header.to-file',
        'punctuation.definition.from-file.diff',
        'punctuation.definition.to-file.diff',
      ],
      settings: {
        foreground: 'var(--code-blue)',
      },
    },
    {
      name: 'Diff Inserted',
      scope: 'markup.inserted.diff',
      settings: {
        foreground: 'var(--code-green)',
      },
    },
    {
      name: 'Diff Deleted',
      scope: 'markup.deleted.diff',
      settings: {
        foreground: 'var(--code-red)',
      },
    },
    {
      name: 'dotenv left-hand side assignments',
      scope: ['variable.other.env'],
      settings: {
        foreground: 'var(--code-blue)',
      },
    },
    {
      name: 'dotenv reference to existing env variable',
      scope: ['string.quoted variable.other.env'],
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      name: 'GDScript functions',
      scope: 'support.function.builtin.gdscript',
      settings: {
        foreground: 'var(--code-blue)',
      },
    },
    {
      name: 'GDScript constants',
      scope: 'constant.language.gdscript',
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      name: 'Comment keywords',
      scope: 'comment meta.annotation.go',
      settings: {
        foreground: 'var(--code-maroon)',
      },
    },
    {
      name: 'go:embed, go:build, etc.',
      scope: 'comment meta.annotation.parameters.go',
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      name: 'Go constants (nil, true, false)',
      scope: 'constant.language.go',
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      name: 'GraphQL variables',
      scope: 'variable.graphql',
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      name: 'GraphQL aliases',
      scope: 'string.unquoted.alias.graphql',
      settings: {
        foreground: 'var(--code-flamingo)',
      },
    },
    {
      name: 'GraphQL enum members',
      scope: 'constant.character.enum.graphql',
      settings: {
        foreground: 'var(--code-teal)',
      },
    },
    {
      name: 'GraphQL field in types',
      scope:
        'meta.objectvalues.graphql constant.object.key.graphql string.unquoted.graphql',
      settings: {
        foreground: 'var(--code-flamingo)',
      },
    },
    {
      name: 'HTML/XML DOCTYPE as keyword',
      scope: [
        'keyword.other.doctype',
        'meta.tag.sgml.doctype punctuation.definition.tag',
        'meta.tag.metadata.doctype entity.name.tag',
        'meta.tag.metadata.doctype punctuation.definition.tag',
      ],
      settings: {
        foreground: 'var(--code-mauve)',
      },
    },
    {
      name: 'HTML/XML-like <tags/>',
      scope: ['entity.name.tag'],
      settings: {
        foreground: 'var(--code-blue)',
        fontStyle: '',
      },
    },
    {
      name: 'Special characters like &amp;',
      scope: [
        'text.html constant.character.entity',
        'text.html constant.character.entity punctuation',
        'constant.character.entity.xml',
        'constant.character.entity.xml punctuation',
        'constant.character.entity.js.jsx',
        'constant.charactger.entity.js.jsx punctuation',
        'constant.character.entity.tsx',
        'constant.character.entity.tsx punctuation',
      ],
      settings: {
        foreground: 'var(--code-red)',
      },
    },
    {
      name: 'HTML/XML tag attribute values',
      scope: ['entity.other.attribute-name'],
      settings: {
        foreground: 'var(--code-yellow)',
      },
    },
    {
      name: 'Components',
      scope: [
        'support.class.component',
        'support.class.component.jsx',
        'support.class.component.tsx',
        'support.class.component.vue',
      ],
      settings: {
        foreground: 'var(--code-pink)',
        fontStyle: '',
      },
    },
    {
      name: 'Annotations',
      scope: ['punctuation.definition.annotation', 'storage.type.annotation'],
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      name: 'Java enums',
      scope: 'constant.other.enum.java',
      settings: {
        foreground: 'var(--code-teal)',
      },
    },
    {
      name: 'Java imports',
      scope: 'storage.modifier.import.java',
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      name: 'Javadoc',
      scope:
        'comment.block.javadoc.java keyword.other.documentation.javadoc.java',
      settings: {
        fontStyle: '',
      },
    },
    {
      name: 'Exported Variable',
      scope: 'meta.export variable.other.readwrite.js',
      settings: {
        foreground: 'var(--code-maroon)',
      },
    },
    {
      name: 'JS/TS constants & properties',
      scope: [
        'variable.other.constant.js',
        'variable.other.constant.ts',
        'variable.other.property.js',
        'variable.other.property.ts',
      ],
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      name: 'JSDoc; these are mainly params, so styled as such',
      scope: [
        'variable.other.jsdoc',
        'comment.block.documentation variable.other',
      ],
      settings: {
        foreground: 'var(--code-maroon)',
        fontStyle: '',
      },
    },
    {
      name: 'JSDoc keywords',
      scope: 'storage.type.class.jsdoc',
      settings: {
        fontStyle: '',
      },
    },
    {
      scope: 'support.type.object.console.js',
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      name: 'Node constants as keywords (module, etc.)',
      scope: ['support.constant.node', 'support.type.object.module.js'],
      settings: {
        foreground: 'var(--code-mauve)',
      },
    },
    {
      name: 'implements as keyword',
      scope: 'storage.modifier.implements',
      settings: {
        foreground: 'var(--code-mauve)',
      },
    },
    {
      name: 'Builtin types',
      scope: [
        'constant.language.null.js',
        'constant.language.null.ts',
        'constant.language.undefined.js',
        'constant.language.undefined.ts',
        'support.type.builtin.ts',
      ],
      settings: {
        foreground: 'var(--code-mauve)',
      },
    },
    {
      scope: 'variable.parameter.generic',
      settings: {
        foreground: 'var(--code-yellow)',
      },
    },
    {
      name: 'Arrow functions',
      scope: [
        'keyword.declaration.function.arrow.js',
        'storage.type.function.arrow.ts',
      ],
      settings: {
        foreground: 'var(--code-teal)',
      },
    },
    {
      name: 'Decorator punctuations (decorators inherit from blue functions, instead of styleguide peach)',
      scope: 'punctuation.decorator.ts',
      settings: {
        foreground: 'var(--code-blue)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Extra JS/TS keywords',
      scope: [
        'keyword.operator.expression.in.js',
        'keyword.operator.expression.in.ts',
        'keyword.operator.expression.infer.ts',
        'keyword.operator.expression.instanceof.js',
        'keyword.operator.expression.instanceof.ts',
        'keyword.operator.expression.is',
        'keyword.operator.expression.keyof.ts',
        'keyword.operator.expression.of.js',
        'keyword.operator.expression.of.ts',
        'keyword.operator.expression.typeof.ts',
      ],
      settings: {
        foreground: 'var(--code-mauve)',
      },
    },
    {
      name: 'Julia macros',
      scope: 'support.function.macro.julia',
      settings: {
        foreground: 'var(--code-teal)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Julia language constants (true, false)',
      scope: 'constant.language.julia',
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      name: 'Julia other constants (these seem to be arguments inside arrays)',
      scope: 'constant.other.symbol.julia',
      settings: {
        foreground: 'var(--code-maroon)',
      },
    },
    {
      name: 'LaTeX preamble',
      scope: 'text.tex keyword.control.preamble',
      settings: {
        foreground: 'var(--code-teal)',
      },
    },
    {
      name: 'LaTeX be functions',
      scope: 'text.tex support.function.be',
      settings: {
        foreground: 'var(--code-sky)',
      },
    },
    {
      name: 'LaTeX math',
      scope: 'constant.other.general.math.tex',
      settings: {
        foreground: 'var(--code-flamingo)',
      },
    },
    {
      name: 'Lua docstring keywords',
      scope:
        'comment.line.double-dash.documentation.lua storage.type.annotation.lua',
      settings: {
        foreground: 'var(--code-mauve)',
        fontStyle: '',
      },
    },
    {
      name: 'Lua docstring variables',
      scope: [
        'comment.line.double-dash.documentation.lua entity.name.variable.lua',
        'comment.line.double-dash.documentation.lua variable.lua',
      ],
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      scope: [
        'heading.1.markdown punctuation.definition.heading.markdown',
        'heading.1.markdown',
        'heading.1.quarto punctuation.definition.heading.quarto',
        'heading.1.quarto',
        'markup.heading.atx.1.mdx',
        'markup.heading.atx.1.mdx punctuation.definition.heading.mdx',
        'markup.heading.setext.1.markdown',
        'markup.heading.heading-0.asciidoc',
      ],
      settings: {
        foreground: 'var(--code-red)',
      },
    },
    {
      scope: [
        'heading.2.markdown punctuation.definition.heading.markdown',
        'heading.2.markdown',
        'heading.2.quarto punctuation.definition.heading.quarto',
        'heading.2.quarto',
        'markup.heading.atx.2.mdx',
        'markup.heading.atx.2.mdx punctuation.definition.heading.mdx',
        'markup.heading.setext.2.markdown',
        'markup.heading.heading-1.asciidoc',
      ],
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      scope: [
        'heading.3.markdown punctuation.definition.heading.markdown',
        'heading.3.markdown',
        'heading.3.quarto punctuation.definition.heading.quarto',
        'heading.3.quarto',
        'markup.heading.atx.3.mdx',
        'markup.heading.atx.3.mdx punctuation.definition.heading.mdx',
        'markup.heading.heading-2.asciidoc',
      ],
      settings: {
        foreground: 'var(--code-yellow)',
      },
    },
    {
      scope: [
        'heading.4.markdown punctuation.definition.heading.markdown',
        'heading.4.markdown',
        'heading.4.quarto punctuation.definition.heading.quarto',
        'heading.4.quarto',
        'markup.heading.atx.4.mdx',
        'markup.heading.atx.4.mdx punctuation.definition.heading.mdx',
        'markup.heading.heading-3.asciidoc',
      ],
      settings: {
        foreground: 'var(--code-green)',
      },
    },
    {
      scope: [
        'heading.5.markdown punctuation.definition.heading.markdown',
        'heading.5.markdown',
        'heading.5.quarto punctuation.definition.heading.quarto',
        'heading.5.quarto',
        'markup.heading.atx.5.mdx',
        'markup.heading.atx.5.mdx punctuation.definition.heading.mdx',
        'markup.heading.heading-4.asciidoc',
      ],
      settings: {
        foreground: 'var(--code-blue)',
      },
    },
    {
      scope: [
        'heading.6.markdown punctuation.definition.heading.markdown',
        'heading.6.markdown',
        'heading.6.quarto punctuation.definition.heading.quarto',
        'heading.6.quarto',
        'markup.heading.atx.6.mdx',
        'markup.heading.atx.6.mdx punctuation.definition.heading.mdx',
        'markup.heading.heading-5.asciidoc',
      ],
      settings: {
        foreground: 'var(--code-mauve)',
      },
    },
    {
      scope: 'markup.bold',
      settings: {
        foreground: 'var(--code-red)',
        fontStyle: 'bold',
      },
    },
    {
      scope: 'markup.italic',
      settings: {
        foreground: 'var(--code-red)',
        fontStyle: 'italic',
      },
    },
    {
      scope: 'markup.strikethrough',
      settings: {
        foreground: 'var(--code-subtext-0)',
        fontStyle: 'strikethrough',
      },
    },
    {
      name: 'Markdown auto links',
      scope: ['punctuation.definition.link', 'markup.underline.link'],
      settings: {
        foreground: 'var(--code-blue)',
      },
    },
    {
      name: 'Markdown links',
      scope: [
        'text.html.markdown punctuation.definition.link.title',
        'text.html.quarto punctuation.definition.link.title',
        'string.other.link.title.markdown',
        'string.other.link.title.quarto',
        'markup.link',
        'punctuation.definition.constant.markdown',
        'punctuation.definition.constant.quarto',
        'constant.other.reference.link.markdown',
        'constant.other.reference.link.quarto',
        'markup.substitution.attribute-reference',
      ],
      settings: {
        foreground: 'var(--code-lavender)',
      },
    },
    {
      name: 'Markdown code spans',
      scope: [
        'punctuation.definition.raw.markdown',
        'punctuation.definition.raw.quarto',
        'markup.inline.raw.string.markdown',
        'markup.inline.raw.string.quarto',
        'markup.raw.block.markdown',
        'markup.raw.block.quarto',
      ],
      settings: {
        foreground: 'var(--code-green)',
      },
    },
    {
      name: 'Markdown triple backtick language identifier',
      scope: 'fenced_code.block.language',
      settings: {
        foreground: 'var(--code-sky)',
      },
    },
    {
      name: 'Markdown triple backticks',
      scope: [
        'markup.fenced_code.block punctuation.definition',
        'markup.raw support.asciidoc',
      ],
      settings: {
        foreground: 'var(--code-overlay-2)',
      },
    },
    {
      name: 'Markdown quotes',
      scope: ['markup.quote', 'punctuation.definition.quote.begin'],
      settings: {
        foreground: 'var(--code-pink)',
      },
    },
    {
      name: 'Markdown separators',
      scope: 'meta.separator.markdown',
      settings: {
        foreground: 'var(--code-teal)',
      },
    },
    {
      name: 'Markdown list bullets',
      scope: [
        'punctuation.definition.list.begin.markdown',
        'punctuation.definition.list.begin.quarto',
        'markup.list.bullet',
      ],
      settings: {
        foreground: 'var(--code-teal)',
      },
    },
    {
      name: 'Quarto headings',
      scope: 'markup.heading.quarto',
      settings: {
        fontStyle: 'bold',
      },
    },
    {
      name: 'Nix attribute names',
      scope: [
        'entity.other.attribute-name.multipart.nix',
        'entity.other.attribute-name.single.nix',
      ],
      settings: {
        foreground: 'var(--code-blue)',
      },
    },
    {
      name: 'Nix parameter names',
      scope: 'variable.parameter.name.nix',
      settings: {
        foreground: 'var(--code-text)',
        fontStyle: '',
      },
    },
    {
      name: 'Nix interpolated parameter names',
      scope: 'meta.embedded variable.parameter.name.nix',
      settings: {
        foreground: 'var(--code-lavender)',
        fontStyle: '',
      },
    },
    {
      name: 'Nix paths',
      scope: 'string.unquoted.path.nix',
      settings: {
        foreground: 'var(--code-pink)',
        fontStyle: '',
      },
    },
    {
      name: 'PHP Attributes',
      scope: ['support.attribute.builtin', 'meta.attribute.php'],
      settings: {
        foreground: 'var(--code-yellow)',
      },
    },
    {
      name: 'PHP Parameters (needed for the leading dollar sign)',
      scope: 'meta.function.parameters.php punctuation.definition.variable.php',
      settings: {
        foreground: 'var(--code-maroon)',
      },
    },
    {
      name: 'PHP Constants (null, __FILE__, etc.)',
      scope: 'constant.language.php',
      settings: {
        foreground: 'var(--code-mauve)',
      },
    },
    {
      name: 'PHP functions',
      scope: 'text.html.php support.function',
      settings: {
        foreground: 'var(--code-sky)',
      },
    },
    {
      name: 'PHPdoc keywords',
      scope: 'keyword.other.phpdoc.php',
      settings: {
        fontStyle: '',
      },
    },
    {
      name: 'Python argument functions reset to text, otherwise they inherit blue from function-call',
      scope: [
        'support.variable.magic.python',
        'meta.function-call.arguments.python',
      ],
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      name: 'Python double underscore functions',
      scope: ['support.function.magic.python'],
      settings: {
        foreground: 'var(--code-sky)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Python `self` keyword',
      scope: [
        'variable.parameter.function.language.special.self.python',
        'variable.language.special.self.python',
      ],
      settings: {
        foreground: 'var(--code-red)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'python keyword flow/logical (for ... in)',
      scope: ['keyword.control.flow.python', 'keyword.operator.logical.python'],
      settings: {
        foreground: 'var(--code-mauve)',
      },
    },
    {
      name: 'python storage type',
      scope: 'storage.type.function.python',
      settings: {
        foreground: 'var(--code-mauve)',
      },
    },
    {
      name: 'python function support',
      scope: [
        'support.token.decorator.python',
        'meta.function.decorator.identifier.python',
      ],
      settings: {
        foreground: 'var(--code-sky)',
      },
    },
    {
      name: 'python function calls',
      scope: ['meta.function-call.python'],
      settings: {
        foreground: 'var(--code-blue)',
      },
    },
    {
      name: 'python function decorators',
      scope: [
        'entity.name.function.decorator.python',
        'punctuation.definition.decorator.python',
      ],
      settings: {
        foreground: 'var(--code-peach)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'python placeholder reset to normal string',
      scope: 'constant.character.format.placeholder.other.python',
      settings: {
        foreground: 'var(--code-pink)',
      },
    },
    {
      name: 'Python exception & builtins such as exit()',
      scope: [
        'support.type.exception.python',
        'support.function.builtin.python',
      ],
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      name: 'entity.name.type',
      scope: ['support.type.python'],
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      name: 'python constants (True/False)',
      scope: 'constant.language.python',
      settings: {
        foreground: 'var(--code-mauve)',
      },
    },
    {
      name: 'Arguments accessed later in the function body',
      scope: ['meta.indexed-name.python', 'meta.item-access.python'],
      settings: {
        foreground: 'var(--code-maroon)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Python f-strings/binary/unicode storage types',
      scope: 'storage.type.string.python',
      settings: {
        foreground: 'var(--code-green)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Python type hints',
      scope: 'meta.function.parameters.python',
      settings: {
        fontStyle: '',
      },
    },
    {
      name: 'Regex string begin/end in JS/TS',
      scope: [
        'string.regexp punctuation.definition.string.begin',
        'string.regexp punctuation.definition.string.end',
      ],
      settings: {
        foreground: 'var(--code-pink)',
      },
    },
    {
      name: 'Regex anchors (^, $)',
      scope: 'keyword.control.anchor.regexp',
      settings: {
        foreground: 'var(--code-mauve)',
      },
    },
    {
      name: 'Regex regular string match',
      scope: 'string.regexp.ts',
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      name: 'Regex group parenthesis & backreference (\\1, \\2, \\3, ...)',
      scope: [
        'punctuation.definition.group.regexp',
        'keyword.other.back-reference.regexp',
      ],
      settings: {
        foreground: 'var(--code-green)',
      },
    },
    {
      name: 'Regex character class []',
      scope: 'punctuation.definition.character-class.regexp',
      settings: {
        foreground: 'var(--code-yellow)',
      },
    },
    {
      name: 'Regex character classes (\\d, \\w, \\s)',
      scope: 'constant.other.character-class.regexp',
      settings: {
        foreground: 'var(--code-pink)',
      },
    },
    {
      name: 'Regex range',
      scope: 'constant.other.character-class.range.regexp',
      settings: {
        foreground: 'var(--code-rosewater)',
      },
    },
    {
      name: 'Regex quantifier',
      scope: 'keyword.operator.quantifier.regexp',
      settings: {
        foreground: 'var(--code-teal)',
      },
    },
    {
      name: 'Regex constant/numeric',
      scope: 'constant.character.numeric.regexp',
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      name: 'Regex lookaheads, negative lookaheads, lookbehinds, negative lookbehinds',
      scope: [
        'punctuation.definition.group.no-capture.regexp',
        'meta.assertion.look-ahead.regexp',
        'meta.assertion.negative-look-ahead.regexp',
      ],
      settings: {
        foreground: 'var(--code-blue)',
      },
    },
    {
      name: 'Rust attribute',
      scope: [
        'meta.annotation.rust',
        'meta.annotation.rust punctuation',
        'meta.attribute.rust',
        'punctuation.definition.attribute.rust',
      ],
      settings: {
        foreground: 'var(--code-yellow)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Rust attribute strings',
      scope: [
        'meta.attribute.rust string.quoted.double.rust',
        'meta.attribute.rust string.quoted.single.char.rust',
      ],
      settings: {
        fontStyle: '',
      },
    },
    {
      name: 'Rust keyword',
      scope: [
        'entity.name.function.macro.rules.rust',
        'storage.type.module.rust',
        'storage.modifier.rust',
        'storage.type.struct.rust',
        'storage.type.enum.rust',
        'storage.type.trait.rust',
        'storage.type.union.rust',
        'storage.type.impl.rust',
        'storage.type.rust',
        'storage.type.function.rust',
        'storage.type.type.rust',
      ],
      settings: {
        foreground: 'var(--code-mauve)',
        fontStyle: '',
      },
    },
    {
      name: 'Rust u/i32, u/i64, etc.',
      scope: 'entity.name.type.numeric.rust',
      settings: {
        foreground: 'var(--code-mauve)',
        fontStyle: '',
      },
    },
    {
      name: 'Rust generic',
      scope: 'meta.generic.rust',
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      name: 'Rust impl',
      scope: 'entity.name.impl.rust',
      settings: {
        foreground: 'var(--code-yellow)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Rust module',
      scope: 'entity.name.module.rust',
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      name: 'Rust trait',
      scope: 'entity.name.trait.rust',
      settings: {
        foreground: 'var(--code-yellow)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Rust struct',
      scope: 'storage.type.source.rust',
      settings: {
        foreground: 'var(--code-yellow)',
      },
    },
    {
      name: 'Rust union',
      scope: 'entity.name.union.rust',
      settings: {
        foreground: 'var(--code-yellow)',
      },
    },
    {
      name: 'Rust enum member',
      scope: 'meta.enum.rust storage.type.source.rust',
      settings: {
        foreground: 'var(--code-teal)',
      },
    },
    {
      name: 'Rust macro',
      scope: [
        'support.macro.rust',
        'meta.macro.rust support.function.rust',
        'entity.name.function.macro.rust',
      ],
      settings: {
        foreground: 'var(--code-blue)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Rust lifetime',
      scope: ['storage.modifier.lifetime.rust', 'entity.name.type.lifetime'],
      settings: {
        foreground: 'var(--code-blue)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Rust string formatting',
      scope: 'string.quoted.double.rust constant.other.placeholder.rust',
      settings: {
        foreground: 'var(--code-pink)',
      },
    },
    {
      name: 'Rust return type generic',
      scope:
        'meta.function.return-type.rust meta.generic.rust storage.type.rust',
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      name: 'Rust functions',
      scope: 'meta.function.call.rust',
      settings: {
        foreground: 'var(--code-blue)',
      },
    },
    {
      name: 'Rust angle brackets',
      scope: 'punctuation.brackets.angle.rust',
      settings: {
        foreground: 'var(--code-sky)',
      },
    },
    {
      name: 'Rust constants',
      scope: 'constant.other.caps.rust',
      settings: {
        foreground: 'var(--code-peach)',
      },
    },
    {
      name: 'Rust function parameters',
      scope: ['meta.function.definition.rust variable.other.rust'],
      settings: {
        foreground: 'var(--code-maroon)',
      },
    },
    {
      name: 'Rust closure variables',
      scope: 'meta.function.call.rust variable.other.rust',
      settings: {
        foreground: 'var(--code-text)',
      },
    },
    {
      name: 'Rust self',
      scope: 'variable.language.self.rust',
      settings: {
        foreground: 'var(--code-red)',
      },
    },
    {
      name: 'Rust metavariable names',
      scope: [
        'variable.other.metavariable.name.rust',
        'meta.macro.metavariable.rust keyword.operator.macro.dollar.rust',
      ],
      settings: {
        foreground: 'var(--code-pink)',
      },
    },
    {
      name: 'Shell shebang',
      scope: [
        'comment.line.shebang',
        'comment.line.shebang punctuation.definition.comment',
        'comment.line.shebang',
        'punctuation.definition.comment.shebang.shell',
        'meta.shebang.shell',
      ],
      settings: {
        foreground: 'var(--code-pink)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Shell shebang command',
      scope: 'comment.line.shebang constant.language',
      settings: {
        foreground: 'var(--code-teal)',
        fontStyle: 'normal',
      },
    },
    {
      name: 'Shell interpolated command',
      scope: [
        'meta.function-call.arguments.shell punctuation.definition.variable.shell',
        'meta.function-call.arguments.shell punctuation.section.interpolation',
        'meta.function-call.arguments.shell punctuation.definition.variable.shell',
        'meta.function-call.arguments.shell punctuation.section.interpolation',
      ],
      settings: {
        foreground: 'var(--code-red)',
      },
    },
    {
      name: 'Shell interpolated command variable',
      scope:
        'meta.string meta.interpolation.parameter.shell variable.other.readwrite',
      settings: {
        foreground: 'var(--code-peach)',
        fontStyle: 'normal',
      },
    },
    {
      scope: [
        'source.shell punctuation.section.interpolation',
        'punctuation.definition.evaluation.backticks.shell',
      ],
      settings: {
        foreground: 'var(--code-teal)',
      },
    },
    {
      name: 'Shell EOF',
      scope: 'entity.name.tag.heredoc.shell',
      settings: {
        foreground: 'var(--code-mauve)',
      },
    },
    {
      name: 'Shell quoted variable',
      scope: 'string.quoted.double.shell variable.other.normal.shell',
      settings: {
        foreground: 'var(--code-text)',
      },
    },
  ],
}

export default theme
