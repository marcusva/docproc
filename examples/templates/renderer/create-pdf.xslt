<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="2.0" 
    xmlns:xsl="http://www.w3.org/1999/XSL/Transform" 
    xmlns:fo="http://www.w3.org/1999/XSL/Format" exclude-result-prefixes="fo">
    <xsl:template match="root">
        <fo:root>
            <fo:layout-master-set>
                <fo:simple-page-master master-name="a4" page-height="29.7cm" page-width="21cm" margin-top="2cm" margin-bottom="2cm" margin-left="2cm" margin-right="2cm">
                    <fo:region-body />
                </fo:simple-page-master>
            </fo:layout-master-set>
            <fo:page-sequence master-reference="a4">
                <fo:flow flow-name="xsl-region-body" font-size="10pt">
                    <fo:block padding-top="3cm">
                        <xsl:value-of select="string-join((customer/firstname, customer/lastname), ' ')"/>
                    </fo:block>
                    <fo:block>
                        <xsl:value-of select="customer/address/street"/>
                    </fo:block>
                    <fo:block>
                        <xsl:value-of select="string-join((customer/address/zip, customer/address/city), ' ')"/>
                    </fo:block>
                    <fo:block font-size="14pt" font-weight="bold" space-before="5cm" margin-bottom="5mm">
                        <xsl:if test="invoice">Your Invoice</xsl:if>
                        <xsl:if test="credit">Your Creidt Note</xsl:if>
                    </fo:block>
                    <xsl:variable name="entry" select="if (exists(invoice)) then invoice else credit"/>
                    <fo:table table-layout="fixed" width="100%" border-collapse="separate">
                        <fo:table-column column-width="4cm"/>
                        <fo:table-column column-width="4cm"/>
                        <fo:table-column column-width="5cm"/>
                        <fo:table-body>
                            <fo:table-row>
                                <fo:table-cell>
                                    <fo:block>
                                        <xsl:value-of select="$entry/date" />
                                    </fo:block>
                                </fo:table-cell>
                                <fo:table-cell>
                                    <fo:block>
                                        <xsl:value-of select="format-number($entry/net, '###,###.00')"/>
                                    </fo:block>
                                </fo:table-cell>
                                <fo:table-cell>
                                    <fo:block>
                                        <xsl:value-of select="format-number($entry/gross, '###,###.00')"/>
                                    </fo:block>
                                </fo:table-cell>
                            </fo:table-row>
                        </fo:table-body>
                    </fo:table>
                </fo:flow>
            </fo:page-sequence>
        </fo:root>
    </xsl:template>
</xsl:stylesheet>